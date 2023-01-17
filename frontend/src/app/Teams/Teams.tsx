import React from 'react'
import { PageSection, PageSectionVariants, TextContent, Text } from '@patternfly/react-core';
import { TeamsWizard } from "@app/Teams/TeamsOnboarding"

export const Teams = () => {

  return (
    <React.Fragment>
      <PageSection style={{
        minHeight: "12%",
        background: "url(https://console.redhat.com/apps/frontend-assets/background-images/new-landing-page/estate_section_banner.svg)",
        backgroundSize: "cover",
        backgroundColor: "black",
        opacity: '0.9'
      }} variant={PageSectionVariants.light}
      >
        <React.Fragment>
          <TextContent style={{ color: "white", display: "inline" }}>
            <div style={{ float: "left", }}>
              <Text component="h2">Get Started with Red Hat Quality Studio</Text>
              <Text component="p">Onboard your team and start to add components to get quality metrics</Text>
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
