package v1

import (
	scheme "github.com/kubermatic/kubermatic/api/pkg/crd/client/clientset/versioned/scheme"
	v1 "github.com/kubermatic/kubermatic/api/pkg/crd/kubermatic/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// AddonsGetter has a method to return a AddonInterface.
// A group's client should implement this interface.
type AddonsGetter interface {
	Addons(namespace string) AddonInterface
}

// AddonInterface has methods to work with Addon resources.
type AddonInterface interface {
	Create(*v1.Addon) (*v1.Addon, error)
	Update(*v1.Addon) (*v1.Addon, error)
	Delete(name string, options *meta_v1.DeleteOptions) error
	DeleteCollection(options *meta_v1.DeleteOptions, listOptions meta_v1.ListOptions) error
	Get(name string, options meta_v1.GetOptions) (*v1.Addon, error)
	List(opts meta_v1.ListOptions) (*v1.AddonList, error)
	Watch(opts meta_v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.Addon, err error)
	AddonExpansion
}

// addons implements AddonInterface
type addons struct {
	client rest.Interface
	ns     string
}

// newAddons returns a Addons
func newAddons(c *KubermaticV1Client, namespace string) *addons {
	return &addons{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the addon, and returns the corresponding addon object, and an error if there is any.
func (c *addons) Get(name string, options meta_v1.GetOptions) (result *v1.Addon, err error) {
	result = &v1.Addon{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("addons").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of Addons that match those selectors.
func (c *addons) List(opts meta_v1.ListOptions) (result *v1.AddonList, err error) {
	result = &v1.AddonList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("addons").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested addons.
func (c *addons) Watch(opts meta_v1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("addons").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

// Create takes the representation of a addon and creates it.  Returns the server's representation of the addon, and an error, if there is any.
func (c *addons) Create(addon *v1.Addon) (result *v1.Addon, err error) {
	result = &v1.Addon{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("addons").
		Body(addon).
		Do().
		Into(result)
	return
}

// Update takes the representation of a addon and updates it. Returns the server's representation of the addon, and an error, if there is any.
func (c *addons) Update(addon *v1.Addon) (result *v1.Addon, err error) {
	result = &v1.Addon{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("addons").
		Name(addon.Name).
		Body(addon).
		Do().
		Into(result)
	return
}

// Delete takes name of the addon and deletes it. Returns an error if one occurs.
func (c *addons) Delete(name string, options *meta_v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("addons").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *addons) DeleteCollection(options *meta_v1.DeleteOptions, listOptions meta_v1.ListOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("addons").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched addon.
func (c *addons) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.Addon, err error) {
	result = &v1.Addon{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("addons").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
