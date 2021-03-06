import { Box, createStyles, Theme, withStyles, WithStyles } from "@material-ui/core";
import { deleteProtectedEndpointAction } from "actions/sso";
import { withSSO, WithSSOProps } from "hoc/withSSO";
import Immutable from "immutable";
import { BasePage } from "pages/BasePage";
import { SSOImplementDetails } from "pages/SSO/Details";
import React from "react";
import { Link } from "react-router-dom";
import {
  ProtectedEndpoint,
  SSOGithubConnector,
  SSOGitlabConnector,
  SSO_CONNECTOR_TYPE,
  SSO_CONNECTOR_TYPE_GITHUB,
  SSO_CONNECTOR_TYPE_GITLAB,
} from "types/sso";
import { CustomizedButton } from "widgets/Button";
import { EditIcon, GithubIcon } from "widgets/Icon";
import { IconButtonWithTooltip } from "widgets/IconButtonWithTooltip";
import { DeleteButtonWithConfirmPopover } from "widgets/IconWithPopover";
import { KPanel } from "widgets/KPanel";
import { KRTable } from "widgets/KRTable";
import { Body, Subtitle1 } from "widgets/Label";
import { KMLink } from "widgets/Link";
import { Loading } from "widgets/Loading";

const styles = (theme: Theme) =>
  createStyles({
    root: {},
  });

interface Props extends WithStyles<typeof styles>, WithSSOProps {}

interface State {}

class SSOPageRaw extends React.PureComponent<Props, State> {
  private renderConnectorDetails = (connector: SSOGitlabConnector | SSOGithubConnector) => {
    // @ts-ignore
    const type = connector.get("type") as SSO_CONNECTOR_TYPE;

    switch (type) {
      case SSO_CONNECTOR_TYPE_GITLAB: {
        const cnt = connector as SSOGitlabConnector;
        const baseURL = cnt.get("config").get("baseURL");
        const groups = cnt.get("config").get("groups") || Immutable.List();
        return (
          <Box key={cnt.get("id")} mt={2}>
            <Subtitle1>
              Gitlab {cnt.get("name")} (
              <KMLink href={baseURL} target="_blank" rel="noopener noreferrer">
                {baseURL}
              </KMLink>
              )
            </Subtitle1>
            Users in groups{" "}
            {groups
              .map((g, index) => (
                <React.Fragment key={index}>
                  <KMLink target="_blank" rel="noopener noreferrer" href={baseURL + "/" + g}>
                    {g}
                  </KMLink>
                  {index < cnt.get("config").get("groups").size - 1 ? ", " : " "}
                </React.Fragment>
              ))
              .toArray()}
          </Box>
        );
      }
      case SSO_CONNECTOR_TYPE_GITHUB: {
        const cnt = connector as SSOGithubConnector;
        return (
          <Box key={cnt.get("id")} mt={2}>
            <Subtitle1>
              <Box display="inline-block" style={{ verticalAlign: "middle" }} mr={1}>
                <GithubIcon />
              </Box>
              Github {cnt.get("name")}
            </Subtitle1>
            {cnt
              .get("config")
              .get("orgs")
              .map((org, index) => {
                const teams = org.get("teams");
                if (teams && teams.size > 0) {
                  return (
                    <Box key={org.get("name")}>
                      Users in organization {org.get("name")} and teams{" "}
                      {org.get("teams").map((team, index) => (
                        <>
                          <a target="_blank" rel="noopener noreferrer" href={"https://github.com/" + team}>
                            {team}
                          </a>
                          {index < teams.size - 1 ? ", " : " "}
                        </>
                      ))}
                    </Box>
                  );
                } else {
                  return <Box key={org.get("name")}>Users in organization {org.get("name")}</Box>;
                }
              })}
          </Box>
        );
      }
    }
  };

  private renderConfigDetails = () => {
    const { ssoConfig } = this.props;

    if (!ssoConfig) {
      return null;
    }

    return (
      <>
        <KPanel title={"Single Sign-on configuration Details"}>
          <Box p={2}>
            <pre>Dex OIDC Issuer: https://{ssoConfig.get("domain")}/dex</pre>
            {ssoConfig.get("connectors") && ssoConfig.get("connectors")!.map(this.renderConnectorDetails)}
          </Box>
          <Box p={2} display="inline-block">
            <CustomizedButton component={Link} size="small" to="/sso/config" variant="outlined" color="primary">
              Edit
            </CustomizedButton>
            {/*{loaded && ssoConfig ? <DangerButton>Delete Single Sign-On Config</DangerButton> : null}*/}
          </Box>
        </KPanel>
        {this.renderProtectedComponents()}
      </>
    );
  };

  private renderNamespace = (rowData: ProtectedEndpoint) => {
    return rowData.get("namespace");
  };

  private renderComponentName = (rowData: ProtectedEndpoint) => {
    return rowData.get("endpointName");
  };

  private renderComponentPorts = (rowData: ProtectedEndpoint) => {
    return !!rowData.get("ports") && rowData.get("ports")!.size > 0 ? rowData.get("ports")!.join(", ") : "All";
  };

  private renderGrantedGroups = (rowData: ProtectedEndpoint) => {
    return !!rowData.get("groups") && rowData.get("groups")!.size > 0 ? rowData.get("groups")!.join(", ") : "All";
  };

  private renderProtectedComponentActions = (rowData: ProtectedEndpoint) => {
    const { dispatch } = this.props;
    return (
      <>
        <IconButtonWithTooltip
          component={Link}
          to={"/sso/endpoints/" + rowData.get("name") + "/edit"}
          size="small"
          tooltipPlacement="top"
          tooltipTitle="Edit"
          aria-label="edit"
        >
          <EditIcon />
        </IconButtonWithTooltip>
        <DeleteButtonWithConfirmPopover
          popupId="delete-sso-popup"
          popupTitle="DELETE SSO?"
          confirmedAction={() => dispatch(deleteProtectedEndpointAction(rowData))}
        />
      </>
    );
  };

  private getKRTableColumns() {
    return [
      {
        Header: "Namespace",
        accessor: "namespace",
      },
      {
        Header: "Component",
        accessor: "component",
      },
      {
        Header: "Ports",
        accessor: "ports",
      },
      {
        Header: "Granted groups",
        accessor: "grantedGroups",
      },
      { Header: "Actions", accessor: "actions" },
    ];
  }

  private getKRTableData() {
    const { protectedEndpoints } = this.props;
    const data: any[] = [];

    protectedEndpoints.forEach((rowData, index) => {
      data.push({
        namespace: this.renderNamespace(rowData),
        component: this.renderComponentName(rowData),
        ports: this.renderComponentPorts(rowData),
        grantedGroups: this.renderGrantedGroups(rowData),
        actions: this.renderProtectedComponentActions(rowData),
      });
    });

    return data;
  }

  private renderKRTable() {
    return <KRTable columns={this.getKRTableColumns()} data={this.getKRTableData()} />;
  }

  private renderProtectedComponents() {
    return (
      <Box mt={2}>
        <KPanel title="Protected Component">
          <Box p={2}>
            <Box display="inline-block" mb={2}>
              <CustomizedButton
                size="small"
                component={Link}
                to="/sso/endpoints/new"
                variant="outlined"
                color="primary"
              >
                New Protected Endpoint
              </CustomizedButton>
            </Box>
            {this.renderKRTable()}
          </Box>
        </KPanel>
      </Box>
    );
  }

  private renderEmptyText = () => {
    return (
      <Body>
        The <strong>single sign-on</strong> feature allows you to configure access permissions for private components.
        Only users with the permissions you configured can access the resources behind. <br />
        Kalm SSO will integrate with your existing user system, such as <strong>github</strong>, <strong>gitlab</strong>
        , <strong>google</strong>, etc.
        <Box mt={2} width={300}>
          <CustomizedButton component={Link} to="/sso/config" variant="contained" color="primary">
            Enable Single Sign-on
          </CustomizedButton>
          {/*{loaded && ssoConfig ? <DangerButton>Delete Single Sign-On Config</DangerButton> : null}*/}
        </Box>
      </Body>
    );
  };

  public render() {
    const { ssoConfig, isSSOConfigLoaded } = this.props;

    if (!isSSOConfigLoaded) {
      return (
        <Box p={2}>
          <Loading />
        </Box>
      );
    }

    return (
      <BasePage>
        <Box p={2}>
          {!!ssoConfig ? this.renderConfigDetails() : this.renderEmptyText()}

          <Box mt={2}>
            <SSOImplementDetails />
          </Box>
        </Box>
      </BasePage>
    );
  }
}

export const SSOPage = withStyles(styles)(withSSO(SSOPageRaw));
