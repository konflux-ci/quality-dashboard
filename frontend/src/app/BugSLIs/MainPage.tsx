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
    ToggleGroup,
    ToggleGroupItem,
} from '@patternfly/react-core';
import { Button } from '@patternfly/react-core';
import { Grid, GridItem } from '@patternfly/react-core';
import { OverviewTable } from './Table';
import { Bug, Info } from './Types';
import { ReactReduxContext, useSelector } from 'react-redux';
import { getBugSLIs, getTeams } from '@app/utils/APIService';
import { validateParam } from '@app/utils/utils';
import { formatDate, getRangeDates } from '@app/Reports/utils';
import { useHistory } from 'react-router-dom';
import { DateTimeRangePicker } from '@app/utils/DateTimeRangePicker';
import { Header } from '@app/utils/Header';
import { InfoBanner } from './InfoBanner';
import { CustomSLIChartDonutCard } from './CustomSLIChartDonutCard';
import { SLIsStackChart } from './SLIsStackChart';

// eslint-disable-next-line prefer-const
export const BugSLIs = () => {
    const [loadingState, setLoadingState] = useState(false);
    const { store } = useContext(ReactReduxContext);
    const state = store.getState();
    const currentTeam = useSelector((state: any) => state.teams.Team);
    const params = new URLSearchParams(window.location.search);
    const [rangeDateTime, setRangeDateTime] = useState(getRangeDates(800));
    const [bugSLIs, setBugSLIs] = useState<Info>();
    const history = useHistory();
    const [isSelected, setIsSelected] = React.useState('resolution');
    const [bugsTable, setBugsTable] = useState<Array<Bug>>([]);

    function handleChange(event, from, to) {
        setRangeDateTime([from, to]);
        params.set('start', formatDate(from));
        params.set('end', formatDate(to));
        history.push(window.location.pathname + '?' + params.toString());
    }

    useEffect(() => {
        if (state.teams.Team != '') {
            getBugSLIs(state.teams.Team, rangeDateTime).then((data: any) => {
                setBugSLIs(data.data)
            });
        }
    }, [rangeDateTime]);

    useEffect(() => {
        if (bugSLIs?.global_sli != undefined) {
            switch (isSelected) {
                case "resolution":
                    setBugsTable(bugSLIs.resolution_time_sli.bugs)
                    setLoadingState(false)
                    break;
                case "response":
                    setLoadingState(false)
                    setBugsTable(bugSLIs.response_time_sli.bugs)
                    break;
                case "triage":
                    setLoadingState(false)
                    setBugsTable(bugSLIs.triage_time_sli.bugs)
                    break;
            }

        }
    }, [isSelected, bugSLIs]);

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
            setBugSLIs({} as Info)

            const team = params.get('team');
            const start = params.get('start');
            const end = params.get('end');

            getBugSLIs(state.teams.Team, rangeDateTime).then((data: any) => {
                if (data.data.length < 1 && (team == state.teams.Team || team == null)) {
                    // setLoadingState(false)
                    history.push('/home/bug-slis?team=' + currentTeam);
                }
                // if (data.data.bugs.length > 0 &&
                if (team == state.teams.Team || team == null) {
                    setBugSLIs(data.data)

                    if (start == null || end == null) {
                        // first click on page or team
                        const start_date = formatDate(rangeDateTime[0]);
                        const end_date = formatDate(rangeDateTime[1]);

                        setLoadingState(false)

                        history.push(
                            '/home/bug-slis?team=' +
                            currentTeam +
                            '&start=' +
                            start_date +
                            '&end=' +
                            end_date
                        );
                    } else {
                        setRangeDateTime([new Date(start), new Date(end)]);

                        history.push(
                            '/home/bug-slis?team=' + currentTeam +
                            '&start=' + start +
                            '&end=' + end
                        );
                    }
                }
            });
        }
    }, [setBugSLIs, currentTeam]);

    const start = rangeDateTime[0];
    const end = rangeDateTime[1];

    const handleItemClick = (isSelected: boolean, event: React.MouseEvent<any> | React.KeyboardEvent | MouseEvent) => {
        const id = event.currentTarget.id;
        setBugsTable([])
        setIsSelected(id);
    };

    return (
        <React.Fragment>
            {/* page title bar */}
            <Header info="Observe which Jira issues are not meeting the defined Bug SLOs."></Header>
            <PageSection variant={PageSectionVariants.light}>
                <Title headingLevel="h3" size={TitleSizes['2xl']}>
                    Bug SLIs
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
                {/* this section will show statistics and details about CiFailures metric */}
                {loadingState && <div style={{ width: '100%', textAlign: "center" }}>
                    <Spinner isSVG diameter="80px" aria-label="Contents of the custom size example" style={{ margin: "100px auto" }} />
                </div>
                }
                {!loadingState && (bugSLIs?.global_sli != undefined) &&
                    (
                        <Grid hasGutter>
                            <GridItem>
                                <DateTimeRangePicker
                                    startDate={start}
                                    endDate={end}
                                    handleChange={(event, from, to) => handleChange(event, from, to)}
                                ></DateTimeRangePicker>
                            </GridItem>
                            <InfoBanner />
                            <GridItem span={3} rowSpan={12}>
                                <CustomSLIChartDonutCard
                                    title="Bug SLI Status"
                                    donutChartColorScale={["#ea2745", "#fbe424", "#61ad50"]}
                                    sli={bugSLIs?.global_sli}
                                    data={[]}
                                    type={""}
                                >
                                </CustomSLIChartDonutCard>

                            </GridItem>
                            <GridItem span={3} rowSpan={12}>
                                <CustomSLIChartDonutCard
                                    title="Resolution Time Bug SLI"
                                    donutChartColorScale={["#ea2745", "#fbe424"]}
                                    sli={{}}
                                    data={bugSLIs?.resolution_time_sli?.bugs}
                                    type={"resolution_sli"}
                                >
                                </CustomSLIChartDonutCard>
                            </GridItem>
                            <GridItem span={3} rowSpan={12}>
                                <CustomSLIChartDonutCard
                                    title="Response Time Bug SLI"
                                    donutChartColorScale={["#ea2745"]}
                                    sli={{}}
                                    data={bugSLIs?.response_time_sli?.bugs}
                                    type={"response_sli"}
                                >
                                </CustomSLIChartDonutCard>
                            </GridItem>
                            <GridItem span={3} rowSpan={12}>
                                <CustomSLIChartDonutCard
                                    title="Triage Time Bug SLI"
                                    donutChartColorScale={["#ea2745", "#fbe424"]}
                                    sli={{}}
                                    data={bugSLIs?.triage_time_sli?.bugs}
                                    type={"triage_sli"}
                                >
                                </CustomSLIChartDonutCard>
                            </GridItem>
                            <GridItem>
                                <Card style={{ textAlign: 'center' }}>
                                    <CardTitle>
                                        Component's Bug SLIs
                                    </CardTitle>
                                    <CardBody>
                                    <SLIsStackChart bugSLIs={bugSLIs}></SLIsStackChart>
                                    </CardBody>
                                </Card>
                            </GridItem>

                            <GridItem>
                                <Card>
                                    <CardTitle>
                                        Bug SLIs Overview
                                    </CardTitle>
                                    <CardBody>
                                        <ToggleGroup aria-label="Default with single selectable">
                                            <ToggleGroupItem
                                                text="Bugs Not Meeting Resolution Time Bug SLO"
                                                buttonId="resolution"
                                                isSelected={isSelected === 'resolution'}
                                                onChange={handleItemClick}
                                            />
                                            <ToggleGroupItem
                                                text="Bugs Not Meeting Response Time Bug SLO"
                                                buttonId="response"
                                                isSelected={isSelected === 'response'}
                                                onChange={handleItemClick}
                                            />
                                            <ToggleGroupItem
                                                text="Bugs Not Meeting Triage Time Bug SLO"
                                                buttonId="triage"
                                                isSelected={isSelected === 'triage'}
                                                onChange={handleItemClick}
                                            />
                                        </ToggleGroup>
                                        <OverviewTable bugSLIs={bugsTable} selected={isSelected}></OverviewTable>
                                    </CardBody>
                                </Card>
                            </GridItem>
                        </Grid>
                    )
                }
            </PageSection>
        </React.Fragment >
    );
};

