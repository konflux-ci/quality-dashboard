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
import { getTeams } from '@app/utils/APIService';

export const Overview = () => {
  const history = useHistory();
  const [noTeams, setNoTeams] = React.useState(false);

  const handleModalToggle = () => {
    setNoTeams(false)
    history.push("/home/teams?isOpen=true")
  };

  const teamsEmpty = () => {
    getTeams().then(res => {
      if (res.data.length == 0) {
        setNoTeams(true)
      }
    })
    if (noTeams) {
      return true
    }

    return false
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
          </Text>
          <Text component="p">
            Observe, track and analyze RHTAP quality metrics.
            By creating a team or joining an existing one, you can be more informed about the code coverage, OpenShift CI prow jobs, and GitHub actions of the RHTAP components.
          </Text>
          <Button onClick={handleModalToggle} type="button" variant="primary"> <PlusIcon></PlusIcon>  {teamsEmpty() ? "Create your first team" : "Create team"} </Button>
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