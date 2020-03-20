package handler

import (
	"encoding/json"
	"github.com/kapp-staging/kapp/api/resources"
	"github.com/kapp-staging/kapp/controller/api/v1alpha1"
	"github.com/labstack/echo/v4"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
)

// Deprecated
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
	applicationList, err := getKappApplicationList(c)

	if err != nil {
		return err
	}

	return c.JSON(200, h.applicationListResponse(c, applicationList))
}

func (h *ApiHandler) handleGetApplicationDetails(c echo.Context) error {
	application, err := getKappApplication(c)

	if err != nil {
		return err
	}

	return c.JSON(200, h.applicationResponse(c, application))
}

// Deprecated
func (h *ApiHandler) handleCreateApplication(c echo.Context) error {
	k8sClient := getK8sClient(c)

	res, err := k8sClient.RESTClient().Post().Body(c.Request().Body).AbsPath("/apis/core.kapp.dev/v1alpha1/namespaces/" + c.Param("namespace") + "/applications").DoRaw()
	if err != nil {
		return err
	}
	return c.JSONBlob(200, res)
}

func (h *ApiHandler) handleCreateApplicationNew(c echo.Context) error {
	application, err := createKappApplication(c)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, h.applicationResponse(c, application))
}

func (h *ApiHandler) handleUpdateApplicationNew(c echo.Context) error {
	application, err := updateKappApplication(c)

	if err != nil {
		return err
	}

	return c.JSON(200, h.applicationResponse(c, application))
}

// Deprecated
func (h *ApiHandler) handleUpdateApplication(c echo.Context) error {
	k8sClient := getK8sClient(c)
	res, err := k8sClient.RESTClient().Put().Body(c.Request().Body).AbsPath("/apis/core.kapp.dev/v1alpha1/namespaces/" + c.Param("namespace") + "/applications/" + c.Param("name")).DoRaw()

	if err != nil {
		return err
	}

	return c.JSONBlob(200, res)
}

func (h *ApiHandler) handleDeleteApplication(c echo.Context) error {
	err := deleteKappApplication(c)
	if err != nil {
		return err
	}
	return c.NoContent(http.StatusNoContent)
}

// Helper functions

func deleteKappApplication(c echo.Context) error {
	k8sClient := getK8sClient(c)
	_, err := k8sClient.RESTClient().Delete().Body(c.Request().Body).AbsPath(kappApplicationUrl(c)).DoRaw()
	return err
}

func createKappApplication(c echo.Context) (*v1alpha1.Application, error) {
	k8sClient := getK8sClient(c)

	crdApplication, err := getApplicationFromContext(c)

	if err != nil {
		return nil, err
	}

	bts, _ := json.Marshal(crdApplication)
	var application v1alpha1.Application
	err = k8sClient.RESTClient().Post().Body(bts).AbsPath(kappApplicationUrl(c)).Do().Into(&application)
	if err != nil {
		return nil, err
	}
	return &application, nil
}

func updateKappApplication(c echo.Context) (*v1alpha1.Application, error) {
	k8sClient := getK8sClient(c)

	crdApplication, err := getApplicationFromContext(c)

	if err != nil {
		return nil, err
	}

	fetched, err := getKappApplication(c)

	if err != nil {
		return nil, err
	}

	crdApplication.ResourceVersion = fetched.ResourceVersion

	bts, _ := json.Marshal(crdApplication)
	var application v1alpha1.Application
	err = k8sClient.RESTClient().Put().Body(bts).AbsPath(kappApplicationUrl(c)).Do().Into(&application)

	if err != nil {
		return nil, err
	}

	return &application, nil
}

func getKappApplication(c echo.Context) (*v1alpha1.Application, error) {
	k8sClient := getK8sClient(c)
	var fetched v1alpha1.Application
	err := k8sClient.RESTClient().Get().AbsPath(kappApplicationUrl(c)).Do().Into(&fetched)
	if err != nil {
		return nil, err
	}
	return &fetched, nil
}

func getKappApplicationList(c echo.Context) (*v1alpha1.ApplicationList, error) {
	k8sClient := getK8sClient(c)
	var fetched v1alpha1.ApplicationList
	err := k8sClient.RESTClient().Get().AbsPath(kappApplicationUrl(c)).Do().Into(&fetched)
	if err != nil {
		return nil, err
	}
	return &fetched, nil
}

func kappApplicationUrl(c echo.Context) string {
	namespace := c.Param("namespace")
	name := c.Param("name")

	if namespace == "" && name == "" {
		return "/apis/core.kapp.dev/v1alpha1/applications"
	}

	if name == "" {
		return "/apis/core.kapp.dev/v1alpha1/namespaces/" + namespace + "/applications"
	}

	return "/apis/core.kapp.dev/v1alpha1/namespaces/" + namespace + "/applications/" + name
}

func getApplicationFromContext(c echo.Context) (*v1alpha1.Application, error) {
	var req resources.CreateOrUpdateApplicationRequest

	if err := c.Bind(&req); err != nil {
		return nil, err
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

	return crdApplication, nil
}

func (h *ApiHandler) applicationResponse(c echo.Context, application *v1alpha1.Application) *resources.ApplicationResponse {
	k8sClient := getK8sClient(c)
	builder := resources.Builder{
		k8sClient,
		h.logger,
	}

	return builder.BuildApplicationDetailsResponse(application)
}

func (h *ApiHandler) applicationListResponse(c echo.Context, applicationList *v1alpha1.ApplicationList) *resources.ApplicationListResponse {
	k8sClient := getK8sClient(c)
	builder := resources.Builder{
		k8sClient,
		h.logger,
	}
	return builder.BuildApplicationListResponse(applicationList)
}