package v1

import (
	v1 "github.com/kubermatic/kubermatic/api/pkg/crd/kubermatic/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// AddonLister helps list Addons.
type AddonLister interface {
	// List lists all Addons in the indexer.
	List(selector labels.Selector) (ret []*v1.Addon, err error)
	// Addons returns an object that can list and get Addons.
	Addons(namespace string) AddonNamespaceLister
	AddonListerExpansion
}

// addonLister implements the AddonLister interface.
type addonLister struct {
	indexer cache.Indexer
}

// NewAddonLister returns a new AddonLister.
func NewAddonLister(indexer cache.Indexer) AddonLister {
	return &addonLister{indexer: indexer}
}

// List lists all Addons in the indexer.
func (s *addonLister) List(selector labels.Selector) (ret []*v1.Addon, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.Addon))
	})
	return ret, err
}

// Addons returns an object that can list and get Addons.
func (s *addonLister) Addons(namespace string) AddonNamespaceLister {
	return addonNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// AddonNamespaceLister helps list and get Addons.
type AddonNamespaceLister interface {
	// List lists all Addons in the indexer for a given namespace.
	List(selector labels.Selector) (ret []*v1.Addon, err error)
	// Get retrieves the Addon from the indexer for a given namespace and name.
	Get(name string) (*v1.Addon, error)
	AddonNamespaceListerExpansion
}

// addonNamespaceLister implements the AddonNamespaceLister
// interface.
type addonNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all Addons in the indexer for a given namespace.
func (s addonNamespaceLister) List(selector labels.Selector) (ret []*v1.Addon, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.Addon))
	})
	return ret, err
}

// Get retrieves the Addon from the indexer for a given namespace and name.
func (s addonNamespaceLister) Get(name string) (*v1.Addon, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1.Resource("addon"), name)
	}
	return obj.(*v1.Addon), nil
}
