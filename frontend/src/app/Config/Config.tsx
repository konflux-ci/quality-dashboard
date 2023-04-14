import React from 'react';
import { Grid, GridItem, PageSection, PageSectionVariants, TextContent, Text } from '@patternfly/react-core';
import { Sample } from './Sample';
import { Editor } from './Editor';


export const Config: React.FunctionComponent = () => (
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
                    <div style={{ float: "left" }}>
                        <Text component="h2">Get Started with Red Hat Quality Studio</Text>
                        <Text component="p">Set the config and start to onboard your teams in a quick way</Text>
                    </div>
                </TextContent>
            </React.Fragment>
        </PageSection>
        <PageSection>
            <Grid hasGutter>
                <GridItem span={8}>
                    <Editor></Editor>
                </GridItem>
                <GridItem span={4}>
                    <Sample></Sample>
                </GridItem>
            </Grid>
        </PageSection>
    </React.Fragment>
);
