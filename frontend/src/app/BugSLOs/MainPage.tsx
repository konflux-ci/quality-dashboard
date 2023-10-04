import React, { useContext, useEffect, useState } from 'react';
import { CopyIcon } from '@patternfly/react-icons';
import {
    PageSection,
    PageSectionVariants,
    Title,
    TitleSizes,
    Spinner,
    Card,
    CardTitle,
    CardBody,
    Flex,
    FlexItem,
} from '@patternfly/react-core';
import { Button } from '@patternfly/react-core';
import { Grid, GridItem } from '@patternfly/react-core';
import { OverviewTable } from './OverviewTable';
import { ReactReduxContext, useSelector } from 'react-redux';
import { getBugSLOsByProject, getTeams } from '@app/utils/APIService';
import { validateParam } from '@app/utils/utils';
import { formatDate, getRangeDates } from '@app/Reports/utils';
import { useHistory } from 'react-router-dom';
import { DateTimeRangePicker } from '@app/utils/DateTimeRangePicker';
import { Header } from '@app/utils/Header';

export interface FailureInfo {
    jira_key: string;
    jira_status: string;
    error_message: string;
    frequency: string;
}

// eslint-disable-next-line prefer-const
export const BugSLOs = () => {
    const [loadingState, setLoadingState] = useState(false);
    const { store } = useContext(ReactReduxContext);
    const state = store.getState();
    const currentTeam = useSelector((state: any) => state.teams.Team);
    const params = new URLSearchParams(window.location.search);
    const [rangeDateTime, setRangeDateTime] = useState(getRangeDates(800));
    const [bugSLOs, setBugSLOs] = useState<any>({});
    const history = useHistory();

    function handleChange(event, from, to) {
        setRangeDateTime([from, to]);
        params.set('start', formatDate(from));
        params.set('end', formatDate(to));
        history.push(window.location.pathname + '?' + params.toString());
    }

    useEffect(() => {
        if (state.teams.Team != '') {
            getBugSLOsByProject(state.teams.Team, rangeDateTime).then((data: any) => {
                setBugSLOs(data.data)
            });
        }
    }, [rangeDateTime]);

    useEffect(() => {
        setLoadingState(true)

        const team = params.get("team")
        if ((team != null) && (team != state.teams.Team)) {
            getTeams().then(res => {
                if (!validateParam(res.data, team)) {
                    setLoadingState(false)
                }
            })
        }

        if (state.teams.Team != '') {
            setBugSLOs([])

            const team = params.get('team');
            const start = params.get('start');
            const end = params.get('end');

            getBugSLOsByProject(state.teams.Team, rangeDateTime).then((data: any) => {
                if (data.data.length < 1 && (team == state.teams.Team || team == null)) {
                    setLoadingState(false)
                    history.push('/home/bug-slos?team=' + currentTeam);
                }

                if (data.data.length > 0 && (team == state.teams.Team || team == null)) {
                    setBugSLOs(data.data)

                    if (start == null || end == null) {
                        // first click on page or team
                        const start_date = formatDate(rangeDateTime[0]);
                        const end_date = formatDate(rangeDateTime[1]);

                        setLoadingState(false)

                        history.push(
                            '/home/bug-slos?team=' +
                            currentTeam +
                            '&start=' +
                            start_date +
                            '&end=' +
                            end_date
                        );
                    } else {
                        setRangeDateTime([new Date(start), new Date(end)]);

                        history.push(
                            '/home/bug-slos?team=' + currentTeam +
                            '&start=' + start +
                            '&end=' + end
                        );
                    }
                    setLoadingState(false)
                }
            });
        }
    }, [setBugSLOs, currentTeam]);


    const start = rangeDateTime[0];
    const end = rangeDateTime[1];


    return (
        <React.Fragment>
            {/* page title bar */}
            <Header info="Observe which JIRA Projects are meeting the defined Bug SLOs."></Header>
            <PageSection variant={PageSectionVariants.light}>
                <Title headingLevel="h3" size={TitleSizes['2xl']}>
                    Bug SLOs
                    <Button
                        onClick={() => navigator.clipboard.writeText(window.location.href)}
                        variant="link"
                        icon={<CopyIcon />}
                        iconPosition="right"
                    >
                        Copy link
                    </Button>
                </Title>
            </PageSection>
            {/* main content  */}
            <PageSection>
                {/* the following toolbar will contain the form (dropdowns and button) to request data to the server */}
                <Grid hasGutter>
                    {/* this section will show statistics and details about CiFailures metric */}
                    {loadingState && <div style={{ width: '100%', textAlign: "center" }}>
                        <Spinner isSVG diameter="80px" aria-label="Contents of the custom size example" style={{ margin: "100px auto" }} />
                    </div>
                    }
                    {!loadingState &&
                        (
                            <GridItem>
                                <Card>
                                    <Flex>
                                        <FlexItem>
                                            <CardTitle>
                                                Bug SLOs Overview
                                            </CardTitle>
                                        </FlexItem>
                                        <FlexItem align={{ default: 'alignRight' }} style={{ marginRight: "25px" }}>
                                            <DateTimeRangePicker
                                                startDate={start}
                                                endDate={end}
                                                handleChange={(event, from, to) => handleChange(event, from, to)}
                                            ></DateTimeRangePicker>
                                        </FlexItem>
                                    </Flex>
                                    <CardBody>
                                        <OverviewTable bugSLOs={bugSLOs}></OverviewTable>
                                    </CardBody>
                                </Card>
                            </GridItem>
                        )
                    }
                </Grid>
            </PageSection>
        </React.Fragment>
    );
};

