import React, { useState, useContext, useEffect } from 'react';
import {
    Card,
    CardTitle,
    CardBody,
    Text,
    PageSection,
    Title,
    Grid,
    GridItem,
    TextContent,
    PageSectionVariants,
    Pagination,
    Chip,
    ChipGroup,
    ToggleGroup,
    ToggleGroupItem,
    Spinner,
} from '@patternfly/react-core';
import {
    TableComposable,
    Thead,
    Tr,
    Th,
    Tbody,
    Td,
    ThProps
} from '@patternfly/react-table';
import { Chart, ChartAxis, ChartGroup, ChartLine, createContainer, ChartThemeColor } from '@patternfly/react-charts';
import { getJirasResolutionTime, getJirasOpen, listE2EBugsKnown } from '@app/utils/APIService';
import { ReactReduxContext, useSelector } from 'react-redux';
import { formatDate, getRangeDates } from '@app/Reports/utils';
import { DateTimeRangePicker } from '@app/utils/DateTimeRangePicker';
import { useHistory } from 'react-router-dom';
import { getLabels } from '@app/utils/utils';
import { help } from '@app/Github/PullRequests';

interface Bugs {
    jira_key: string;
    created_at: string;
    deleted_at: string;
    updated_at: string;
    resolved_at: string;
    resolution_time: string;
    last_change_time: string;
    status: string;
    summary: string;
    affects_versions: string;
    fix_versions: string;
    components: string;
    labels: string;
    url: string;
    teams_bugs: string;
}

export const Jira = () => {
    const { store } = useContext(ReactReduxContext);
    const state = store.getState();
    const dispatch = store.dispatch;
    const currentTeam = useSelector((state: any) => state.teams.Team);
    const [bugsKnown, setBugsKnown] = useState<any>({});
    const BugsAffectingCI = "Bugs Affecting CI"
    const [rangeDateTime, setRangeDateTime] = useState(getRangeDates(30));
    const priorities = ["Global", "Major", "Critical", "Blocker", "Normal", "Undefined", "Minor"]
    const [loadingState, setLoadingState] = useState(false);

    const getJiraData = (ID) => {
        setLoadingState(true)
        const newData = {}

        if (currentTeam != "") {
            const promises = new Array()
            priorities.forEach(p => {
                promises.push(getJirasOpen(p, currentTeam, rangeDateTime))
                promises.push(getJirasResolutionTime(p, currentTeam, rangeDateTime))
            })

            Promise.all(promises).then(function (values) {
                values.map(value => {
                    if (value.data.hasOwnProperty("open")) {
                        if (!newData.hasOwnProperty(value.data.open.priority)) {
                            newData[value.data.open.priority] = {}
                        }
                        newData[value.data.open.priority].open = value.data.open
                    } else if (value.data.hasOwnProperty("resolution_time")) {
                        if (!newData.hasOwnProperty(value.data.resolution_time.priority)) {
                            newData[value.data.open.priority] = {}
                        }
                        newData[value.data.resolution_time.priority].resolved = value.data.resolution_time
                    }
                })
                setApiDataCache(newData)
                setSelected(ID)
            });
        }
    }

    useEffect(() => {
        if (currentTeam != "") {
            listE2EBugsKnown().then(res => {
                const bugs = new Array<Bugs>
                res.data.forEach((bug, _) => {
                    bugs.push({
                        jira_key: bug.key,
                        created_at: bug.fields.created,
                        deleted_at: "",
                        updated_at: bug.fields.updated,
                        resolved_at: "",
                        resolution_time: "",
                        last_change_time: "",
                        status: bug.fields.status.description,
                        summary: bug.fields.summary,
                        affects_versions: "",
                        fix_versions: "",
                        components: "",
                        labels: "",
                        url: "https://issues.redhat.com/browse/" + bug.key,
                        teams_bugs: "",
                    });
                })
                setBugsKnown(bugs)
            })

            const selected = params.get("selected")
            const ID =  selected != null ? selected : "Global"
            getJiraData(ID)

            const start = params.get("start")
            const end = params.get("end")

            if (start == null || end == null) {
                history.push('/home/jira?team=' + currentTeam + '&selected=' + ID + '&start=' + formatDate(rangeDateTime[0]) + '&end=' + formatDate(rangeDateTime[1]))
            } else {
                setRangeDateTime([new Date(start), new Date(end)])
                history.push('/home/jira?team=' + currentTeam + '&selected=' + ID + '&start=' + start + '&end=' + end)
            }
        }
    }, [currentTeam]);

    const [selected, setSelected] = useState<string>('');
    const [apiDataCache, setApiDataCache] = useState<any>({});
    const [resolutionTimeChart, setResolutionTimeChart] = useState<any>({});
    const [bugsChart, setBugsChart] = useState<any>({});
    const [bugsTable, setBugsTable] = useState<any>({});
    const [graphicsVisible, setGraphicsVisible] = useState(false);
    // longVersionVisible indicates if 'resolved_at' and 'resolution_time' should be displayed in bugs table
    const [longVersionVisible, setLongVersionVisible] = useState(false);
    const history = useHistory();
    const params = new URLSearchParams(window.location.search);

    const [isSelected, setIsSelected] = React.useState('open');

    const handleItemClick = (isSelected: boolean, event: React.MouseEvent<any> | React.KeyboardEvent | MouseEvent) => {
        const id = event.currentTarget.id;
        setBugsTable([])
        setIsSelected(id);
    };

    useEffect(() => {
        const selected = params.get("selected")
        const target = selected != null ? selected : "Global"
        getJiraData(target)
    }, [rangeDateTime]);

    const onClick = (event: React.MouseEvent) => {
        const ID = event.currentTarget.id
        if (selected != ID) {
            if (!apiDataCache.hasOwnProperty(ID)) {
                const promise0 = getJirasOpen(ID, state.teams.Team, rangeDateTime)
                const promise1 = getJirasResolutionTime(ID, state.teams.Team, rangeDateTime)
                Promise.all([promise0, promise1]).then(function (values) {
                    let newData = {}
                    newData[ID] = {}
                    newData[ID].resolved = values[1].data.resolution_time
                    newData[ID].open = values[0].data.open
                    setApiDataCache({
                        ...apiDataCache,
                        ...newData
                    })
                });
            }
        }
        setSelected(ID)
        params.set("selected", ID)
        history.push(window.location.pathname + '?' + params.toString());
    };

    useEffect(() => {
        if (apiDataCache[selected] && selected != BugsAffectingCI) {
            let rtc = new Array()
            let bc = new Array()
            let rbt = new Array()
            let obt = new Array()
            let obc = new Array()

            apiDataCache[selected].resolved.days.map((item, index) => {
                rtc.push({
                    name: "Resolution Time (" + selected + ")",
                    x: new Date(item.name).toLocaleDateString('en-US', { day: 'numeric', month: 'short', year: '2-digit' }),
                    y: item.total
                })
                bc.push({
                    name: "Resolved Bugs (" + selected + ")",
                    x: new Date(item.name).toLocaleDateString('en-US', { day: 'numeric', month: 'short', year: '2-digit' }),
                    y: item.resolved_bugs
                })
                rbt = [...rbt, ...item.bugs]
            })
            apiDataCache[selected].open.days.map((item, index) => {
                let date = item.name.match(/([^_]+)/g)
                obc.push({
                    name: "Open Bugs (" + selected + ")",
                    x: new Date(item.name).toLocaleDateString('en-US', { day: 'numeric', month: 'short', year: '2-digit' }),
                    y: item.open_bugs
                })
                obt = [...obt, ...item.bugs]
            })

            setBugsChart([bc, obc])
            setResolutionTimeChart([rtc])
            setGraphicsVisible(true)
            if (isSelected == 'resolved') {
                setBugsTable(rbt)
                setLongVersionVisible(true)
            }
            if (isSelected == 'open') {
                setBugsTable(obt)
                setLongVersionVisible(false)
            }
        }
        if (selected == BugsAffectingCI) {
            setLongVersionVisible(false)
            setGraphicsVisible(false)
            setBugsChart([])
            setResolutionTimeChart([])
            setBugsTable(bugsKnown)
        }
        setLoadingState(false)
    }, [selected, isSelected, apiDataCache]);

    function onBarChartClick(event) {
        console.log("clicked", event)
    }

    function handleChange(event, from, to) {
        setRangeDateTime([from, to])
        params.set("start", formatDate(from))
        params.set("end", formatDate(to))
        history.push(window.location.pathname + '?' + params.toString());
    }

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
                <TextContent style={{ color: "white", display: "inline" }}>
                    <div style={{ float: "left", }}>
                        <Text component="h2">Get Started with Red Hat Quality Studio</Text>
                        <Text component="p">Observe which Jira Issues are affecting the CI pass rate.</Text>
                    </div>
                </TextContent>
            </PageSection>
            <PageSection>
                {loadingState && <div style={{ width: '100%', textAlign: "center" }}>
                    <Spinner isSVG diameter="80px" aria-label="Contents of the custom size example" style={{ margin: "100px auto" }} />
                </div>
                }
                {!loadingState && <React.Fragment>
                    <Grid hasGutter>
                        <GridItem>
                            <DateTimeRangePicker
                                startDate={rangeDateTime[0]}
                                endDate={rangeDateTime[1]}
                                handleChange={(event, from, to) => handleChange(event, from, to)}
                            >
                            </DateTimeRangePicker>
                        </GridItem>
                        {graphicsVisible && <GridItem order={{ default: "2" }}>
                            <Grid hasGutter sm={6} md={6} lg={6} xl={6}>
                                <GridItem order={{ default: "1" }}>
                                    <Card style={{ textAlign: 'center' }}>
                                        <CardTitle style={{ textAlign: 'center' }}>Average Resolution Time</CardTitle>
                                        <CardBody>
                                            <Title headingLevel='h1' size="2xl">
                                                {apiDataCache[selected] &&
                                                    <span>
                                                        <span>{parseFloat(apiDataCache[selected].resolved.total) || "-"}</span>
                                                        <span style={{ paddingLeft: '5px', fontSize: '15px', fontWeight: 'normal' }}>day(s)</span>
                                                    </span>
                                                }
                                                {!apiDataCache[selected] && "-"}
                                            </Title>
                                            <BugsChart chartTitle="Average Resolution Time" data={resolutionTimeChart} onBarClick={onBarChartClick}></BugsChart>
                                        </CardBody>
                                    </Card>
                                </GridItem>
                                <GridItem order={{ default: "2" }}>
                                    <Card style={{ textAlign: 'center' }}>
                                        <CardTitle>Bugs</CardTitle>
                                        <CardBody>
                                            <Title headingLevel='h1' size="2xl">
                                                {apiDataCache[selected] &&
                                                    <span>
                                                        <span>{apiDataCache[selected].open.open_bugs || "-"}</span>
                                                        <span style={{ fontSize: '15px', paddingRight: '10px' }}> open </span>
                                                        <span>{apiDataCache[selected].resolved.resolved_bugs || "-"}</span>
                                                        <span style={{ fontSize: '15px', paddingRight: '10px' }}> resolved </span>
                                                    </span>
                                                }
                                                {!apiDataCache[selected] && "-"}
                                            </Title>
                                            <BugsChart chartTitle="Bugs" data={bugsChart} onBarClick={onBarChartClick}></BugsChart>
                                        </CardBody>
                                    </Card>
                                </GridItem>
                            </Grid>
                        </GridItem>}
                        <GridItem order={{ default: "1" }}>
                            <Grid hasGutter span={3}>
                                <GridItem order={{ default: "1" }}>
                                    <Card isSelectable onClick={onClick} style={{ textAlign: 'center' }} isSelected={selected.includes('Global')} id="Global">
                                        <CardTitle>Global</CardTitle>
                                        <CardBody>
                                            <Title headingLevel='h1' size="2xl">
                                                {apiDataCache["Global"] ? <span>{apiDataCache["Global"].open.open_bugs} <span style={{ fontSize: '10px' }}>open</span></span> : "-"}
                                                <span style={{ marginLeft: '20px' }}>&nbsp;</span>
                                                {apiDataCache["Global"] ? <span>{apiDataCache["Global"].resolved.resolved_bugs} <span style={{ fontSize: '10px' }}>resolved</span></span> : "-"}
                                            </Title>
                                        </CardBody>
                                    </Card>
                                </GridItem>
                                <GridItem order={{ default: "2" }}>
                                    <Card isSelectable onClick={onClick} style={{ textAlign: 'center' }} isSelected={selected.includes('Blocker')} id="Blocker">
                                        <CardTitle>Blockers</CardTitle>
                                        <CardBody>
                                            <Title headingLevel='h1' size="2xl">
                                                {apiDataCache["Blocker"] ? <span>{apiDataCache["Blocker"].open.open_bugs} <span style={{ fontSize: '10px' }}>open</span></span> : "-"}
                                                <span style={{ marginLeft: '20px' }}>&nbsp;</span>
                                                {apiDataCache["Blocker"] ? <span>{apiDataCache["Blocker"].resolved.resolved_bugs} <span style={{ fontSize: '10px' }}>resolved</span></span> : "-"}
                                            </Title>
                                        </CardBody>
                                    </Card>
                                </GridItem>
                                <GridItem order={{ default: "3" }}>
                                    <Card isSelectable onClick={onClick} style={{ textAlign: 'center' }} isSelected={selected.includes('Critical')} id="Critical">
                                        <CardTitle>Critical Bugs</CardTitle>
                                        <CardBody>
                                            <Title headingLevel='h1' size="2xl">
                                                {apiDataCache["Critical"] ? <span>{apiDataCache["Critical"].open.open_bugs} <span style={{ fontSize: '10px' }}>open</span></span> : "-"}
                                                <span style={{ marginLeft: '20px' }}>&nbsp;</span>
                                                {apiDataCache["Critical"] ? <span>{apiDataCache["Critical"].resolved.resolved_bugs} <span style={{ fontSize: '10px' }}>resolved</span></span> : "-"}
                                            </Title>
                                        </CardBody>
                                    </Card>
                                </GridItem>
                                <GridItem order={{ default: "4" }}>
                                    <Card isSelectable style={{ textAlign: 'center' }} onClick={onClick} isSelected={selected.includes('Major')} id="Major">
                                        <CardTitle>Major Bugs</CardTitle>
                                        <CardBody>
                                            <Title headingLevel='h1' size="2xl">
                                                {apiDataCache["Major"] ? <span>{apiDataCache["Major"].open.open_bugs} <span style={{ fontSize: '10px' }}>open</span></span> : "-"}
                                                <span style={{ marginLeft: '20px' }}>&nbsp;</span>
                                                {apiDataCache["Major"] ? <span>{apiDataCache["Major"].resolved.resolved_bugs} <span style={{ fontSize: '10px' }}>resolved</span></span> : "-"}
                                            </Title>
                                        </CardBody>
                                    </Card>
                                </GridItem>
                                <GridItem order={{ default: "5" }}>
                                    <Card isSelectable onClick={onClick} style={{ textAlign: 'center' }} isSelected={selected.includes('Normal')} id="Normal">
                                        <CardTitle>Normal Bugs</CardTitle>
                                        <CardBody>
                                            <Title headingLevel='h1' size="2xl">
                                                {apiDataCache["Normal"] ? <span>{apiDataCache["Normal"].open.open_bugs} <span style={{ fontSize: '10px' }}>open</span></span> : "-"}
                                                <span style={{ marginLeft: '20px' }}>&nbsp;</span>
                                                {apiDataCache["Normal"] ? <span>{apiDataCache["Normal"].resolved.resolved_bugs} <span style={{ fontSize: '10px' }}>resolved</span></span> : "-"}
                                            </Title>
                                        </CardBody>
                                    </Card>
                                </GridItem>
                                <GridItem order={{ default: "6" }}>
                                    <Card isSelectable style={{ textAlign: 'center' }} onClick={onClick} isSelected={selected.includes('Minor')} id="Minor">
                                        <CardTitle>Minor Bugs</CardTitle>
                                        <CardBody>
                                            <Title headingLevel='h1' size="2xl">
                                                {apiDataCache["Minor"] ? <span>{apiDataCache["Minor"].open.open_bugs} <span style={{ fontSize: '10px' }}>open</span></span> : "-"}
                                                <span style={{ marginLeft: '20px' }}>&nbsp;</span>
                                                {apiDataCache["Minor"] ? <span>{apiDataCache["Minor"].resolved.resolved_bugs} <span style={{ fontSize: '10px' }}>resolved</span></span> : "-"}
                                            </Title>
                                        </CardBody>
                                    </Card>
                                </GridItem>
                                <GridItem order={{ default: "7" }}>
                                    <Card isSelectable style={{ textAlign: 'center' }} onClick={onClick} isSelected={selected.includes('Undefined')} id="Undefined">
                                        <CardTitle>Undefined Bugs</CardTitle>
                                        <CardBody>
                                            <Title headingLevel='h1' size="2xl">
                                                {apiDataCache["Undefined"] ? <span>{apiDataCache["Undefined"].open.open_bugs} <span style={{ fontSize: '10px' }}>open</span></span> : "-"}
                                                <span style={{ marginLeft: '20px' }}>&nbsp;</span>
                                                {apiDataCache["Undefined"] ? <span>{apiDataCache["Undefined"].resolved.resolved_bugs} <span style={{ fontSize: '10px' }}>resolved</span></span> : "-"}
                                            </Title>
                                        </CardBody>
                                    </Card>
                                </GridItem>
                                <GridItem order={{ default: "8" }}>
                                    <Card isSelectable onClick={onClick} style={{ textAlign: 'center' }} isSelected={selected.includes(BugsAffectingCI)} id={BugsAffectingCI}>
                                        <CardTitle>
                                            Bugs affecting CI
                                            {help("Bugs affecting CI in the projects DEVHAS, SRVKP, GITOPSRVCE, HACBS, RHTAP, and RHTAPBUGS.")}
                                        </CardTitle>
                                        <CardBody>
                                            <Title headingLevel='h1' size="2xl">
                                                {bugsKnown ? <span>{bugsKnown.length} <span style={{ fontSize: '10px' }}>open</span></span> : "-"}
                                            </Title>
                                        </CardBody>
                                    </Card>
                                </GridItem>
                            </Grid>
                        </GridItem>
                        <GridItem order={{ default: "3" }}>
                            <Card style={{ fontSize: "12px" }}>
                                <CardTitle>Bugs</CardTitle>
                                <CardBody>
                                    <Grid hasGutter span={2}>
                                        {(selected != BugsAffectingCI) && <GridItem order={{ default: "1" }}>
                                            <ToggleGroup aria-label="Default with single selectable">
                                                <ToggleGroupItem
                                                    text="Open bugs"
                                                    buttonId="open"
                                                    isSelected={isSelected === 'open'}
                                                    onChange={handleItemClick}
                                                />
                                                <ToggleGroupItem
                                                    text="Resolved bugs"
                                                    buttonId="resolved"
                                                    isSelected={isSelected === 'resolved'}
                                                    onChange={handleItemClick}
                                                />
                                            </ToggleGroup>
                                        </GridItem>}
                                        <GridItem order={{ default: "2" }}>
                                            <ChipGroup categoryName="Active filters: " numChips={5}>
                                                <Chip key={selected} isReadOnly style={{ fontSize: '15px' }}>
                                                    {selected} bugs
                                                </Chip>
                                            </ChipGroup>
                                        </GridItem>
                                    </Grid>
                                    <ComposableTableStripedTr bugs={bugsTable} longVersion={longVersionVisible}></ComposableTableStripedTr>
                                </CardBody>
                            </Card>
                        </GridItem>
                    </Grid>
                </React.Fragment>}
            </PageSection>

        </React.Fragment>
    )
}
const BugsChart: React.FC<{ chartTitle: string, data: any, onBarClick: any }> = ({ chartTitle, data, onBarClick }) => {
    const ZoomVoronoiContainer = createContainer("zoom", "voronoi");

    let legendData: { name: string }[] = []
    if (data.length > 0) {
        legendData = data.map((dataset, index) => {
            return { name: dataset[0]["name"] }
        })
    }

    const getMaxY = (data) => {
        let maxY = 0;
        data.map((dataset, idx) => {
            dataset.forEach((data) => {
                if (data.y > maxY) {
                    maxY = data.y;
                }
            })
        })

        return maxY
    };

    const maxY = getMaxY(data);

    return (
        <div style={{ margin: '0 auto', height: '60%', width: '90%', marginTop: '15px' }}>
            {data.length > 0 &&
                <Chart
                    height={210}
                    containerComponent={
                        <ZoomVoronoiContainer
                            labels={({ datum }) => chartTitle == "Bugs" ? getLabels(datum, "Resolved Bugs") : getLabels(datum, "Resolution Time")}
                            voronoiDimension="x"
                            voronoiPadding={0}
                            constrainToVisibleArea
                        />
                    }
                    legendData={legendData}
                    legendPosition='bottom'
                    padding={{
                        bottom: 70,
                        left: 40,
                        right: 40,
                        top: 20
                    }}
                    name="chart"
                    themeColor={ChartThemeColor.blue}
                    maxDomain={{ y: maxY + 1 }}
                    minDomain={{ y: 0 }}
                >

                    <ChartAxis fixLabelOverlap={true} style={{ axisLabel: { fontSize: 8, padding: 30 }, tickLabels: { fontSize: 7 } }} />
                    <ChartAxis fixLabelOverlap={true} dependentAxis={true} showGrid style={{ axisLabel: { fontSize: 8, padding: 30 }, tickLabels: { fontSize: 8 } }} />
                    {data.length > 0 &&
                        <ChartGroup>
                            {data.map((dataset, index) => (
                                <ChartLine
                                    key={index}
                                    style={{
                                        data: { strokeWidth: 2 },
                                        parent: { border: "1px solid #ccc" },
                                        labels: { fontSize: '10px' }
                                    }}
                                    data={dataset}
                                />
                            ))}
                        </ChartGroup>
                    }
                </Chart>
            }
        </div>
    );
}

const ComposableTableStripedTr: React.FC<{ bugs: any, longVersion: boolean }> = ({ bugs, longVersion }) => {
    const [bugsPage, setBugsPage] = useState<Array<Bugs>>([]);
    const [page, setPage] = React.useState(1);
    const [perPage, setPerPage] = React.useState(10);
    const [count, setCount] = useState(bugs.length);
    const [filters, setFilters] = useState({});
    const [activeSortIndex, setActiveSortIndex] = React.useState<number | null>(null);
    const [activeSortDirection, setActiveSortDirection] = React.useState<'asc' | 'desc' | null>(null);

    useEffect(() => {
        if (bugs.length == 0) {
            setPage(1)
            setBugsPage([])
        }
        if (bugs.length > 0) {
            setBugsPage(bugs.slice(0, perPage))
            setPage(1)
        }
    }, [bugs]);

    const columnNames = {
        jira_key: "ID",
        created_at: "Created at",
        deleted_at: "Deleted at",
        updated_at: "Updated at",
        resolved_at: "Resolved at",
        resolution_time: "Resolution time",
        last_change_time: "Last changed at",
        status: "Status",
        summary: "Summary",
        affects_versions: "Affected versions",
        fix_versions: "Fix versions",
        components: "Components",
        labels: "Labels",
        url: "URL",
        teams_bugs: "Team"
    };

    useEffect(
        () => {
            setCount(bugs.length);
        },
        [bugs],
    );

    useEffect(() => {
        setCount(bugs.length)
        if (bugs.length > 0) {
            const filteredRows = filterRows(bugs, filters)
            const sortedRows = sortRows(filteredRows)
            setCount(sortedRows.length)

            const from = (page - 1) * perPage
            const to = (page - 1) * perPage + perPage >= sortedRows.length ? sortedRows.length : (page - 1) * perPage + perPage;

            setBugsPage(sortedRows.slice(from, to))
        }
    }, [page, perPage, filters, activeSortIndex, activeSortDirection]);


    const onSetPage = (_event: React.MouseEvent | React.KeyboardEvent | MouseEvent, newPage: number) => {
        setPage(newPage);
    };

    const onPerPageSelect = (
        _event: React.MouseEvent | React.KeyboardEvent | MouseEvent,
        newPerPage: number,
        newPage: number
    ) => {
        setPerPage(newPerPage);
        setPage(newPage);
    };


    // Filters helpers
    const columns = [
        { column: 'jira_key', label: 'ID' },
        { column: 'summary', label: 'Summary' },
        { column: 'status', label: 'Status' },
        { column: 'created_at', label: 'Created at' },
        { column: 'updated_at', label: 'Updated at' },
    ]

    if (longVersion) {
        columns.push({ column: 'resolved_at', label: 'Resolved at' })
        columns.push({ column: 'resolution_time', label: 'Resolution Time' })
    }

    function filterRows(rows, filters) {
        if (Object.keys(filters).length === 0) return rows

        return rows.filter(row => {
            return Object.keys(filters).every(column => {
                const value = row[column]
                const searchValue = filters[column]

                // handle Resolution Time filter
                if (typeof value === "number") {
                    return value.toFixed(2) + "h" == searchValue
                }

                // handle ID, Summary, Created at, Updated at, and Resolved at filters
                if (typeof value === 'string') {
                    return value.toLocaleLowerCase().includes(searchValue.toLocaleLowerCase())
                }
                return false
            })
        })
    }

    const handleSearch = (value, column) => {
        if (value) {
            setFilters(prevFilters => ({
                ...prevFilters,
                [column]: value,
            }))
        } else {
            setFilters(prevFilters => {
                const updatedFilters = { ...prevFilters }
                delete updatedFilters[column]

                return updatedFilters
            })
        }
    }
    // End of filter helpers


    // Sort helpers
    const getSortableRowValues = (bug: Bugs): (string | number)[] => {
        const { jira_key, summary, status, created_at, updated_at, resolved_at, resolution_time } = bug;
        return [jira_key, summary, status, created_at, updated_at, resolved_at, resolution_time];
    };

    const sortRows = (rows) => {
        if (activeSortIndex !== null) {
            return rows.sort((a, b) => {
                const aValue = getSortableRowValues(a)[activeSortIndex] ? getSortableRowValues(a)[activeSortIndex] : "-";
                const bValue = getSortableRowValues(b)[activeSortIndex] ? getSortableRowValues(b)[activeSortIndex] : "-";
                if (typeof aValue === 'number') {
                    // Numeric sort
                    if (activeSortDirection === 'asc') {
                        return (aValue as number) - (bValue as number);
                    }
                    return (bValue as number) - (aValue as number);
                } else {
                    // String sort
                    if (activeSortDirection === 'asc') {
                        return (aValue as string).localeCompare(bValue as string, undefined, { numeric: true, sensitivity: 'base' });
                    }
                    return (bValue as string).localeCompare(aValue as string, undefined, { numeric: true, sensitivity: 'base' });
                }
            });
        }
        return rows
    }

    const getSortParams = (columnIndex: number): ThProps['sort'] => ({
        sortBy: {
            index: activeSortIndex as number,
            direction: activeSortDirection as any
        },
        onSort: (_event, index, direction) => {
            setActiveSortIndex(index);
            setActiveSortDirection(direction);
        },
        columnIndex
    });
    // End of sort helpers

    return (
        <div>
            <Pagination
                perPageComponent="button"
                itemCount={count}
                perPage={perPage}
                page={page}
                onSetPage={onSetPage}
                widgetId="top-example"
                onPerPageSelect={onPerPageSelect}
            />

            <TableComposable aria-label="Simple table" >
                <Thead>
                    <Tr>
                        {columns.map((column, idx) => {
                            return (
                                <Th width={10} sort={getSortParams(idx)} key={idx}>
                                    {column.label}
                                </Th>
                            )
                        })}
                    </Tr>
                    <Tr>
                        {columns.map(c => {
                            return (
                                <Th key={c.column}>
                                    <input style={{ width: '100%' }}
                                        key={`${c.column}-search`}
                                        type="search"
                                        placeholder={`Search`}
                                        value={filters[c.column]}
                                        onChange={event => handleSearch(event.target.value, c.column)}
                                    />
                                </Th>
                            )
                        })}
                    </Tr>
                </Thead>
                <Tbody>
                    {bugsPage.map((bug, index) => (
                        <Tr key={bug.jira_key} {...(index % 2 === 0 && { isStriped: true })}>
                            <Td dataLabel={columnNames.jira_key}><a href={bug.url} target={bug.url}>{bug.jira_key}</a></Td>
                            <Td dataLabel={columnNames.summary}>{bug.summary}</Td>
                            <Td dataLabel={columnNames.status}>{bug.status ? bug.status : "-"}</Td>
                            <Td dataLabel={columnNames.created_at}>{formatDate(new Date(bug.created_at))}</Td>
                            <Td dataLabel={columnNames.updated_at}>{formatDate(new Date(bug.updated_at))}</Td>
                            {longVersion && <Td dataLabel={columnNames.resolved_at}>{formatDate(new Date(bug.resolved_at))}</Td>}
                            {longVersion && <Td dataLabel={columnNames.resolution_time}>{!Number.isNaN(parseFloat(bug.resolution_time)) ? parseFloat(bug.resolution_time) + " day(s)" : "-"}</Td>}
                        </Tr>
                    ))}
                </Tbody>
            </TableComposable>

            <Pagination
                perPageComponent="button"
                itemCount={count}
                perPage={perPage}
                page={page}
                onSetPage={onSetPage}
                widgetId="top-example"
                onPerPageSelect={onPerPageSelect}
            />
        </div>
    );
};
