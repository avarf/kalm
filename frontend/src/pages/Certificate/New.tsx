import React from "react";
import { createStyles, Theme, withStyles, WithStyles, Grid } from "@material-ui/core";
import { connect } from "react-redux";
import { TDispatchProp } from "types";
import { CertificateFormType, newEmptyCertificateForm } from "types/certificate";
import { createCertificateAction } from "actions/certificate";
import { CertificateForm } from "forms/Certificate";
import { BasePage } from "pages/BasePage";
import { H6 } from "widgets/Label";
import { push } from "connected-react-router";

const styles = (theme: Theme) =>
  createStyles({
    root: {},
  });

export interface Props extends WithStyles<typeof styles>, TDispatchProp {}

class CertificateNewRaw extends React.PureComponent<Props> {
  private submit = async (certificate: CertificateFormType) => {
    try {
      const { dispatch } = this.props;
      await dispatch(createCertificateAction(certificate, false));
    } catch (e) {
      console.log(e);
    }
  };

  private onSubmitSuccess = () => {
    this.props.dispatch(push("/certificates"));
  };

  public render() {
    const { classes } = this.props;
    return (
      <BasePage secondHeaderRight={<H6>New Certificate</H6>}>
        <div className={classes.root}>
          <Grid container spacing={2}>
            <Grid item xs={8} sm={8} md={8}>
              <CertificateForm
                onSubmitSuccess={this.onSubmitSuccess}
                onSubmit={this.submit}
                initialValues={newEmptyCertificateForm}
              />
            </Grid>
          </Grid>
        </div>
      </BasePage>
    );
  }
}

export const CertificateNewPage = withStyles(styles)(connect()(CertificateNewRaw));
