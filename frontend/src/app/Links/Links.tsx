/* eslint-disable @typescript-eslint/explicit-module-boundary-types */
/* eslint-disable react/jsx-key */
import React from 'react';
import {
  PageSection,
  TextContent,
  Text,
  PageSectionVariants,
} from '@patternfly/react-core';
import { StatusPages } from './StatusPages';

export const Links = () => {
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
            Important links for Red Hat Trusted Application CI
          </Text>
          <Text component="p">
            This page contains important links for RHTAP CI 3rd party services
          </Text>
        </TextContent>
      </PageSection>
      <PageSection>
        <StatusPages />
      </PageSection>
    </React.Fragment>
  );
}