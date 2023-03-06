import React from "react";
import { Alert, Card, CardBody, CardTitle } from "@patternfly/react-core";
import { ExternalLinkAltIcon } from "@patternfly/react-icons";


const getColor = (coveredFixed) => {
    if (coveredFixed >= 0 && coveredFixed <= 33.33) {
        return "red"

    } else if (coveredFixed >= 33.33 && coveredFixed <= 66.66) {
        return "orange"
    }
    return "green"
}

export const CodeCov = (props) => {

    const repo = props.repo
    const coveredFixed = repo.coverage.coverage_percentage

    return (
        <Card style={{ width: '100%', height: '100%', textAlign: 'center' }}>
            <CardTitle>
                <div>
                    CodeCov
                </div>
                <div style={{ color: "grey", fontSize: 12 }}>
                    {repo ? (
                        <a href={`https://app.codecov.io/gh/${repo.organization}/${repo.repoName}`} target="blank" rel="noopener noreferrer">More info</a>
                    ) : "N/A"}
                </div>
            </CardTitle>
           <CardBody>
                <div style={{ fontSize: 25 }}>
                    {repo ? (
                        <div style={{ color: getColor(coveredFixed) }}>{coveredFixed.toFixed(2) + "%"}</div>
                    ) : "N/A"}
                </div>
            </CardBody>
        </Card>
    );
}