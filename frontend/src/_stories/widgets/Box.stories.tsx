import { action } from "@storybook/addon-actions";
import { select } from "@storybook/addon-knobs";
import Immutable from "immutable";
import React from "react";
import { ImmutableMap } from "typings";
import { ErrorBadge, PendingBadge, SuccessBadge } from "widgets/Badge";
import { FlexRowItemCenterBox } from "widgets/Box";
import { DeleteIcon, EditIcon } from "widgets/Icon";
import { IconButtonWithTooltip } from "widgets/IconButtonWithTooltip";
import { KRTable } from "widgets/KRTable";

export default {
  title: "Widgets/Box",
  component: FlexRowItemCenterBox,
};

export const FlexRowItemCenterBoxExample = () => (
  <FlexRowItemCenterBox style={{ border: "1px solid red" }}>
    <ErrorBadge /> 2
  </FlexRowItemCenterBox>
);

interface CertificateContent {
  name: string;
  isSelfManaged: boolean;
  selfManagedCertContent: string;
  selfManagedCertPrivateKey: string;
  httpsCertIssuer: string;
  domains: Immutable.List<string>;
  ready: string;
  reason: string;
}

type Certificate = ImmutableMap<CertificateContent>;

const renderStatus = (rowData: Certificate) => {
  const ready = rowData.get("ready");

  if (ready === "True") {
    return (
      <FlexRowItemCenterBox>
        <FlexRowItemCenterBox mr={1} style={{ border: "1px solid red" }}>
          <SuccessBadge />
        </FlexRowItemCenterBox>
        <FlexRowItemCenterBox style={{ border: "1px solid red" }}>Normal</FlexRowItemCenterBox>
      </FlexRowItemCenterBox>
    );
  } else if (!!rowData.get("reason")) {
    return (
      <FlexRowItemCenterBox style={{ border: "1px solid red" }}>
        <FlexRowItemCenterBox mr={1} style={{ border: "1px solid red" }}>
          <PendingBadge />
        </FlexRowItemCenterBox>
        <FlexRowItemCenterBox style={{ border: "1px solid red" }}>{rowData.get("reason")}</FlexRowItemCenterBox>
      </FlexRowItemCenterBox>
    );
  } else {
    return <PendingBadge />;
  }
};

const renderName = (rowData: Certificate) => {
  return rowData.get("name");
};

const renderDomains = (rowData: Certificate) => {
  return (
    <>
      {rowData.get("domains")?.map((domain) => {
        return <div key={domain}>{domain}</div>;
      })}
    </>
  );
};

const renderMoreActions = (rowData: Certificate) => {
  return (
    <>
      {rowData.get("isSelfManaged") && (
        <IconButtonWithTooltip tooltipTitle="Edit" aria-label="edit" onClick={action("Edit")}>
          <EditIcon />
        </IconButtonWithTooltip>
      )}
      <IconButtonWithTooltip tooltipTitle="Delete" aria-label="delete" onClick={action("Delete")}>
        <DeleteIcon />
      </IconButtonWithTooltip>
    </>
  );
};

const renderType = (rowData: Certificate) => {
  return rowData.get("isSelfManaged") ? "UPLOADED" : "MANAGED";
};

const getCertifiate = () => {
  const domains = Immutable.List(["kalm.dev"]);
  return Immutable.Map({
    name: "Certificates for Domains",
    domains: domains,
    isSelfManaged: true,
    ready: select("Ready status", ["True", "False"], "True", "Certificate"),
    reason: select("Reason", ["Pending Connect...", "Waiting Provider"], "Pending", "Certificate"),
  }) as Certificate;
};

const getKRTableColumns = () => {
  return [
    {
      Header: "Name",
      accessor: "name",
    },
    {
      Header: "Domains",
      accessor: "domains",
    },
    {
      Header: "Status",
      accessor: "status",
    },
    {
      Header: "Type",
      accessor: "isSelfManaged",
    },
    {
      Header: "Actions",
      accessor: "moreAction",
    },
  ];
};

const getKRTableData = () => {
  const data: any[] = [];

  const rowData = getCertifiate() as Certificate;

  data.push({
    name: renderName(rowData),
    domains: renderDomains(rowData),
    status: renderStatus(rowData),
    isSelfManaged: renderType(rowData),
    moreAction: renderMoreActions(rowData),
  });

  return data;
};

export const FlexRowItemCenterBoxWithTable = () => {
  return <KRTable columns={getKRTableColumns()} data={getKRTableData()} />;
};
