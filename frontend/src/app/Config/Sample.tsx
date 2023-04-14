import React from 'react';
import { CodeBlock, CodeBlockAction, CodeBlockCode, ClipboardCopyButton, Card, CardTitle } from '@patternfly/react-core';

export const Sample: React.FunctionComponent = () => {
    const [copied, setCopied] = React.useState(false);

    const clipboardCopyFunc = (event, text) => {
        navigator.clipboard.writeText(text.toString());
    };

    const onClick = (event, text) => {
        clipboardCopyFunc(event, text);
        setCopied(true);
    };

    const code = `teams:
   - name: team-example
     description: description-example
     jira_projects:
        - STONE
     repositories:
        - name: e2e-tests
          organization: redhat-appstudio
        - name: quality-dashboard
          organization: redhat-appstudio
`;

    const actions = (
        <React.Fragment>
            <CodeBlockAction>
                <ClipboardCopyButton
                    id="basic-copy-button"
                    textId="code-content"
                    aria-label="Copy to clipboard"
                    onClick={e => onClick(e, code)}
                    exitDelay={copied ? 1500 : 600}
                    maxWidth="110px"
                    variant="plain"
                >
                Copy to clipboard
                </ClipboardCopyButton>
            </CodeBlockAction>
        </React.Fragment>
    );

    return (
        <Card style={{ backgroundColor: "#f1f7fc" }}>
        <CardTitle>Sample</CardTitle>
            <CodeBlock actions={actions} style={{ backgroundColor: "#f1f7fc", color: "grey" }}>
            <CodeBlockCode id="code-content">{code}</CodeBlockCode>
        </CodeBlock>
        </Card>
    );
};
