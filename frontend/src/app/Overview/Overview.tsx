/* eslint-disable @typescript-eslint/explicit-module-boundary-types */
/* eslint-disable react/jsx-key */
import React from 'react';
import {
  PageSection,
  TextContent,
  Text,
  PageSectionVariants,
  Button
} from '@patternfly/react-core';
import { PlusIcon } from '@patternfly/react-icons';
import { InfoBanner } from './InfoBanner';
import { About } from './About';
import { useHistory } from 'react-router-dom';

export const Overview = () => {
  const history = useHistory();

  const handleModalToggle = () => {
    history.push("/home/teams?isOpen=true")
  };

  return (
    <React.Fragment>
      <PageSection style={{
        minHeight: "12%",
        background: "url(https://console.redhat.com/apps/frontend-assets/background-images/new-landing-page/estate_section_banner.svg)",
        backgroundSize: "cover",
        backgroundColor: "black",
        opacity: '0.9'
      }} variant={PageSectionVariants.light}>
        <TextContent style={{ color: "white" }}>
          <Text component="h2">
            Red Hat Quality Studio
          </Text>
          <Text component="p">
            Observe, track and analyze StoneSoup quality metrics.
            By creating a team or joining an existing one, you can be more informed about the code coverage, OpenShift CI prow jobs, and GitHub actions of the StoneSoup components.
          </Text>
          <Button onClick={handleModalToggle} type="button" variant="primary"> <PlusIcon></PlusIcon> Create Team </Button>
        </TextContent>
      </PageSection>
      <PageSection>
        <InfoBanner />
      </PageSection>
      <PageSection isFilled>
        <About />
      </PageSection>
    </React.Fragment>
  );
}