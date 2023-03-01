import React from "react";
import { Alert, Card, CardBody, CardTitle } from "@patternfly/react-core";
import { ExternalLinkAltIcon } from "@patternfly/react-icons";


const rederCoverageEffects = (repo) => {
    const coveredFixed = repo.coverage.coverage_percentage
    if (coveredFixed >= 0 && coveredFixed <= 33.33) {
        return <Alert title={coveredFixed.toFixed(2) + "%"} variant="danger" isInline isPlain />
    } else if (coveredFixed >= 33.33 && coveredFixed <= 66.66) {
        return <Alert title={coveredFixed.toFixed(2) + "%"} variant="warning" isInline isPlain />
    }
    return <Alert title={coveredFixed.toFixed(2) + "%"} variant="success" isInline isPlain />
}

export const CodeCov = (props) => {

    const repo = props.repo

    return (
        <Card style={{ width: "100%", height: "100%" }}>
            <CardTitle>
                {repo ? (
                    <a href={`https://app.codecov.io/gh/${repo.organization}/${repo.repoName}`}>CodeCov<ExternalLinkAltIcon style={{ marginLeft: "0.5%" }}></ExternalLinkAltIcon></a>
                ) : "CodeCov"}
            </CardTitle>
            <CardBody>
                {repo ? (
                    rederCoverageEffects(repo)
                ) : "N/A"}
            </CardBody>
        </Card>
    );
}