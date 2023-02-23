/* eslint-disable react/jsx-key */
import React, { useState } from 'react';
import {
  PageSection,
  TextContent,
  Text,
  PageSectionVariants,
  Button
} from '@patternfly/react-core';
import { CopyIcon, PlusIcon } from '@patternfly/react-icons';
import { InfoBanner } from './InfoBanner';
import { About } from './About';
import { TeamForm } from '@app/Teams/TeamForm';

export const Overview = () => {
  const [isOpen, setOpen] = useState<boolean>(false);

  const handleModalToggle = () => {
    setOpen(true)
  };

  function handleChange(event, open) {
    setOpen(open)
  }

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
            <Button onClick={() => navigator.clipboard.writeText(window.location.href)} variant="link" icon={<CopyIcon />} iconPosition="right">
              Copy link
            </Button>
          </Text>
          <Text component="p">
            Observe, track and analyze StoneSoup quality metrics.
            By creating a team or joining an existing one, you can be more informed about the code coverage, OpenShift CI prow jobs, and GitHub actions of the StoneSoup components.
          </Text>
          <Button onClick={handleModalToggle} type="button" variant="primary"> <PlusIcon></PlusIcon> Add Team </Button>
          <TeamForm
            isOpen={isOpen}
            handleChange={(event, open) => handleChange(event, open)}
          ></TeamForm>
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