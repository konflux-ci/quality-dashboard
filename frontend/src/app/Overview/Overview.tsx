/* eslint-disable react/jsx-key */
import React from 'react';
import {
  PageSection,
  TextContent,
  Text,
  PageSectionVariants,
  DrawerContentBody,
  Button
} from '@patternfly/react-core';
import { RepositoriesTable } from '@app/Repositories/RepositoriesTable';
import { CopyIcon } from '@patternfly/react-icons';
import { isValidTeam } from '@app/utils/utils';
import { InfoBanner } from './InfoBanner';
import { About } from './About';

export const Overview = () => {

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
        </TextContent>
      </PageSection>
      <PageSection>
        <InfoBanner />
      </PageSection>
      <PageSection isFilled>
        <About />
      </PageSection>
      {isValidTeam() && <PageSection style={{
        minHeight: "12%"
      }}>
        <RepositoriesTable showTableToolbar={true} showCoverage={true} showDescription={false} enableFiltersOnTheseColumns={['git_organization']}></RepositoriesTable>
        <React.Fragment>
        </React.Fragment>
      </PageSection>
      }
      <PageSection padding={{ default: 'noPadding' }}>
        <DrawerContentBody hasPadding></DrawerContentBody>
      </PageSection>
    </React.Fragment>
  );
}