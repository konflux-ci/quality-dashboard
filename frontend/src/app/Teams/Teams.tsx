import React, { useContext } from 'react'
import {PageSection, PageSectionVariants, TextContent, Text, Button, ButtonVariant } from '@patternfly/react-core';
import { Context } from '@app/store/store';
import { TeamsWizard } from "@app/Teams/TeamsOnboarding"

export const Teams = () => {

  return (
    <React.Fragment>
      <PageSection style={{
          minHeight : "12%",
          background:"url(https://console.redhat.com/apps/frontend-assets/background-images/new-landing-page/estate_section_banner.svg)",
          backgroundSize: "cover",
          backgroundColor : "black",
          opacity: '0.9'
        }} variant={PageSectionVariants.light}
      >
        <React.Fragment>
          <TextContent style={{color: "white", display:"inline"}}>
            <div style={{float: "left", }}>
              <Text component="h2">Get Started with Red Hat Quality Studio</Text>
              <Text component="p">Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.</Text>
            </div>
          </TextContent>
        </React.Fragment>
      </PageSection>
      <PageSection>
        <TeamsWizard></TeamsWizard>
      </PageSection>
    </React.Fragment>
    )
}
