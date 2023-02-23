import React from 'react'
import { PageSection, PageSectionVariants, TextContent, Text, DrawerContentBody } from '@patternfly/react-core';
import { RepositoriesTable } from './RepositoriesTable';

export const Repositories = () => {

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
                            <Text component="p">Manage your team's repositories and check the code coverage.</Text>
                        </div>
                    </TextContent>
                </React.Fragment>
            </PageSection>
            <PageSection style={{
                minHeight: "12%"
            }}>
                <RepositoriesTable showTableToolbar={true} showCoverage={true} showDescription={false} enableFiltersOnTheseColumns={['git_organization']}></RepositoriesTable>
                <React.Fragment>
                </React.Fragment>
            </PageSection>
            <PageSection padding={{ default: 'noPadding' }}>
                <DrawerContentBody hasPadding></DrawerContentBody>
            </PageSection>
        </React.Fragment>
    )
}
