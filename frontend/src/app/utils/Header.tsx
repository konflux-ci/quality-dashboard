import React from 'react';
import { TextContent, Text, PageSectionVariants, PageSection } from '@patternfly/react-core';

export const Header = (props) => {
    const info = props?.info

    return (
        <PageSection style={{
            minHeight: "12%",
            background: "url(https://console.redhat.com/apps/frontend-assets/background-images/new-landing-page/estate_section_banner.svg)",
            backgroundSize: "cover",
            backgroundColor: "black",
            opacity: '0.9'
        }} variant={PageSectionVariants.light}
        >
            <TextContent style={{ color: "white", display: "inline" }}>
                <div style={{ float: "left", }}>
                    <Text component="h2">Get Started with Red Hat Quality Studio's Plugins</Text>
                    <Text component="p">{info}</Text>
                </div>
            </TextContent>

        </PageSection>
    )
}
