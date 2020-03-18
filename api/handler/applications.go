package handler

import (
	"encoding/json"
	"github.com/kapp-staging/kapp/api/resources"
	"github.com/kapp-staging/kapp/controller/api/v1alpha1"
	"github.com/labstack/echo/v4"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (h *ApiHandler) handleGetApplicationsOld(c echo.Context) error {
	k8sClient := getK8sClient(c)
	namespace := c.Param("namespace")

	var path string

	if namespace == "" {
		path = "/apis/core.kapp.dev/v1alpha1/applications"
	} else {
		path = "/apis/core.kapp.dev/v1alpha1/namespaces/" + namespace + "/applications"
	}

	bts, err := k8sClient.RESTClient().Get().AbsPath(path).DoRaw()

	if err != nil {
		return err
	}

	return c.JSONBlob(200, bts)
}

func (h *ApiHandler) handleGetApplications(c echo.Context) error {
	k8sClient := getK8sClient(c)
	namespace := c.Param("namespace")

	var path string

	if namespace == "" {
		path = "/apis/core.kapp.dev/v1alpha1/applications"
	} else {
		path = "/apis/core.kapp.dev/v1alpha1/namespaces/" + namespace + "/applications"
	}

	var applicationList v1alpha1.ApplicationList

	err := k8sClient.RESTClient().Get().AbsPath(path).Do().Into(&applicationList)

	if err != nil {
		return err
	}

	builder := resources.Builder{
		k8sClient,
		h.logger,
	}

	return c.JSON(200, builder.BuildApplicationListResponse(&applicationList))
}

func (h *ApiHandler) handleGetApplicationDetails(c echo.Context) error {
	k8sClient := getK8sClient(c)
	namespace := c.Param("namespace")
	name := c.Param("name")
	var application v1alpha1.Application
	err := k8sClient.RESTClient().Get().AbsPath("/apis/core.kapp.dev/v1alpha1/namespaces/" + namespace + "/applications/" + name).Do().Into(&application)

	if err != nil {
		return err
	}

	builder := resources.Builder{
		k8sClient,
		h.logger,
	}

	return c.JSON(200, builder.BuildApplicationDetailsResponse(&application))
}

func (h *ApiHandler) handleCreateApplication(c echo.Context) error {
	k8sClient := getK8sClient(c)

	res, err := k8sClient.RESTClient().Post().Body(c.Request().Body).AbsPath("/apis/core.kapp.dev/v1alpha1/namespaces/" + c.Param("namespace") + "/applications").DoRaw()
	if err != nil {
		return err
	}
	return c.JSONBlob(200, res)
}

func (h *ApiHandler) handleCreateApplicationNew(c echo.Context) error {
	k8sClient := getK8sClient(c)

	var req resources.CreateApplicationRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	crdApplication := &v1alpha1.Application{
		TypeMeta: metaV1.TypeMeta{
			Kind:       "Application",
			APIVersion: "core.kapp.dev/v1alpha1",
		},
		ObjectMeta: metaV1.ObjectMeta{
			Namespace: req.Application.Namespace,
			Name:      req.Application.Name,
		},
		Spec: v1alpha1.ApplicationSpec{
			SharedEnv:  req.Application.SharedEnvs,
			Components: req.Application.Components,
		},
	}

	bts, _ := json.Marshal(crdApplication)
	var application v1alpha1.Application
	err := k8sClient.RESTClient().Post().Body(bts).AbsPath("/apis/core.kapp.dev/v1alpha1/namespaces/" + c.Param("namespace") + "/applications").Do().Into(&application)

	if err != nil {
		return err
	}

	builder := resources.Builder{
		k8sClient,
		h.logger,
	}

	return c.JSON(200, builder.BuildApplicationDetailsResponse(&application))
}

func (h *ApiHandler) handleUpdateApplication(c echo.Context) error {
	k8sClient := getK8sClient(c)
	res, err := k8sClient.RESTClient().Put().Body(c.Request().Body).AbsPath("/apis/core.kapp.dev/v1alpha1/namespaces/" + c.Param("namespace") + "/applications/" + c.Param("name")).DoRaw()

	if err != nil {
		return err
	}

	return c.JSONBlob(200, res)
}

func (h *ApiHandler) handleDeleteApplication(c echo.Context) error {
	k8sClient := getK8sClient(c)
	res, err := k8sClient.RESTClient().Delete().AbsPath("/apis/core.kapp.dev/v1alpha1/namespaces/" + c.Param("namespace") + "/applications/" + c.Param("name")).DoRaw()

	if err != nil {
		return err
	}

	return c.JSONBlob(200, res)
}
