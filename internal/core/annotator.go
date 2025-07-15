package core

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/alizmhdi/shield/internal/config"
	"github.com/alizmhdi/shield/internal/k8s"
	v1 "k8s.io/api/networking/v1"
)

const WhitelistAnnotation = "nginx.ingress.kubernetes.io/whitelist-source-range"

type Annotator struct {
	K8s    *k8s.Client
	Config *config.Config
}

func NewAnnotator(k8sClient *k8s.Client, cfg *config.Config) *Annotator {
	return &Annotator{
		K8s:    k8sClient,
		Config: cfg,
	}
}

// ApplyAnnotations applies the whitelist annotation to the configured Ingress resources.
func (a *Annotator) ApplyAnnotations() error {
	ctx := context.Background()
	for _, assign := range a.Config.IngressAssignments {
		source, ok := a.Config.Whitelists[assign.Whitelist]
		if !ok {
			return fmt.Errorf("whitelist '%s' not found in config", assign.Whitelist)
		}

		ips := make([]string, 0)
		if len(source.IPs) > 0 {
			ips = append(ips, source.IPs...)
		}
		if source.URL != "" {
			remoteIPs, err := a.fetchIPsFromURL(source.URL)
			if err != nil {
				return fmt.Errorf("failed to fetch IPs from %s: %w", source.URL, err)
			}
			ips = append(ips, remoteIPs...)
		}
		if len(ips) == 0 {
			return fmt.Errorf("no IPs found for whitelist '%s'", assign.Whitelist)
		}
		ipList := a.joinIPs(ips)

		if assign.Name != "" && assign.Namespace != "" {
			ingresses, err := a.K8s.ListIngresses(ctx, assign.Namespace)
			if err != nil {
				return err
			}
			for _, ing := range ingresses {
				if ing.Name == assign.Name {
					a.setAnnotation(&ing, ipList)
					err := a.K8s.UpdateIngress(ctx, assign.Namespace, &ing)
					if err != nil {
						return err
					}
					fmt.Printf("Annotated %s/%s\n", assign.Namespace, assign.Name)
				}
			}
		} else if assign.Namespace != "" {
			ingresses, err := a.K8s.ListIngresses(ctx, assign.Namespace)
			if err != nil {
				return err
			}
			for _, ing := range ingresses {
				a.setAnnotation(&ing, ipList)
				err := a.K8s.UpdateIngress(ctx, assign.Namespace, &ing)
				if err != nil {
					return err
				}
				fmt.Printf("Annotated %s/%s\n", assign.Namespace, ing.Name)
			}
		} else {
			// TODO: Apply to all ingresses in all namespaces if needed
			fmt.Fprintf(os.Stderr, "Skipping assignment with no namespace: %+v\n", assign)
		}
	}
	return nil
}

// fetchIPsFromURL retrieves a list of IPs from a remote JSON URL.
func (a *Annotator) fetchIPsFromURL(url string) ([]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var ips []string
	if err := json.Unmarshal(body, &ips); err != nil {
		return nil, err
	}
	return ips, nil
}

// setAnnotation sets the whitelist annotation on the given Ingress.
func (a *Annotator) setAnnotation(ing *v1.Ingress, ipList string) {
	if ing.Annotations == nil {
		ing.Annotations = make(map[string]string)
	}
	ing.Annotations[WhitelistAnnotation] = ipList
}

// joinIPs joins a slice of IPs into a comma-separated string.
func (a *Annotator) joinIPs(ips []string) string {
	result := ""
	for i, ip := range ips {
		if i > 0 {
			result += ","
		}
		result += ip
	}
	return result
}
