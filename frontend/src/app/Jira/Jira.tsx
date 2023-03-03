import React, { useContext, useEffect, useRef, useLayoutEffect, useState } from 'react';
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
    Tooltip
} from '@patternfly/react-core';
import {
    TableComposable,
    Thead,
    Tr,
    Th,
    Tbody,
    Td,
    Caption
} from '@patternfly/react-table';
import { Chart, ChartAxis, ChartGroup, ChartLine,ChartBar, createContainer } from '@patternfly/react-charts';
import { getJirasResolutionTime } from '@app/utils/APIService';

interface Bugs {
    jira_key: string;
    created_at: string;
    deleted_at: string;
    updated_at: string;
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

function getMonth(monthStr){
    return new Date(monthStr+'-1-01').getMonth()+1
}

export const Jira = () => {

    useEffect(() => {
        const ID = "Global"
        getJirasResolutionTime("").then((res) => {
            if (res.code === 200) {
                const result = res.data;

                // Cache results from API
                let newData = {}
                newData[ID] = result
                setApiDataCache({
                    ...apiDataCache,
                    ...newData
                })
                setSelected(ID)
            }
            
        })
    }, []);

    const [selected, setSelected] = useState<string>('');
    const [apiDataCache, setApiDataCache] = useState<any>({});
    const [resolutionTimeChart, setResolutionTimeChart] = useState<any>({});
    const [bugsChart, setBugsChart] = useState<any>({});
    const [bugsTable, setBugsTable] = useState<any>({});
    const [stats, setStats] = useState<any>({});

    const [isSelected, setIsSelected] = React.useState('resolved');
    const handleItemClick = (isSelected: boolean, event: React.MouseEvent<any> | React.KeyboardEvent | MouseEvent) => {
        const id = event.currentTarget.id;
        setIsSelected(id);
    };

    const onClick = (event: React.MouseEvent) => {
        let ID = event.currentTarget.id
        if(selected != ID){
            if(!apiDataCache.hasOwnProperty(ID)){
                getJirasResolutionTime(ID).then((res) => {
                    if (res.code === 200) {
                        const result = res.data;

                        // Cache results from API
                        let newData = {}
                        newData[ID] = result
                        setApiDataCache({
                            ...apiDataCache,
                            ...newData
                        })
                        setSelected(ID)
                    }
                    
                })
            }
            setSelected(ID)
        } else {
            setSelected(ID)
        }
    };
    
    useEffect(() => {
        if(apiDataCache[selected]){
            let rtc = new Array(12)
            let bc = new Array(12)
            let bt = new Array()
            let obc = new Array(12)

            apiDataCache[selected].resolution_time.months.map(item => {
                rtc[getMonth(item.name)-1] = { 
                    name: "Resolution Time", 
                    x: item.name.slice(0, 3), 
                    y: item.total
                }
                bc[getMonth(item.name)-1] = { 
                    name: "Resolved Bugs", 
                    x: item.name.slice(0, 3), 
                    y: item.resolved_bugs
                }
                obc[getMonth(item.name)-1] = { 
                    name: "Opened Bugs", 
                    x: item.name.slice(0, 3), 
                    y: item.resolved_bugs + Math.floor(Math.random() * 25)
                }
                bt = [...bt, ...item.bugs]
            })
            setBugsChart([bc, obc])
            setResolutionTimeChart([rtc])
            setBugsTable(bt)
        }
    }, [selected, apiDataCache]);

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
                        <Text component="p">Observe which Jira Issues are affecting CI pass rate.</Text>
                    </div>
                </TextContent>
            </PageSection>
            <PageSection>
                <React.Fragment>
                <Grid hasGutter>
                    <GridItem order={{default: "2"}}>
                        <Grid hasGutter sm={6} md={6} lg={6} xl={6}>
                            <GridItem order={{default: "1"}}>
                                <Card id="tooltip-boundary">
                                    <CardTitle>Average Resolution Time</CardTitle>
                                    <CardBody>
                                        <Title headingLevel='h1' size="4xl">
                                            { apiDataCache[selected] &&
                                            <span>
                                                <span>{ parseFloat(apiDataCache[selected].resolution_time.total).toFixed(2) || "-"}</span>
                                                <span style={{paddingLeft: '5px', fontSize: '15px', fontWeight: 'normal'}}>hours</span>
                                            </span>
                                            }
                                            { !apiDataCache[selected] && "-" }
                                        </Title>
                                        <BugsChart chartType="line" data={resolutionTimeChart}></BugsChart>
                                    </CardBody>
                                </Card>
                            </GridItem>
                            <GridItem order={{default: "2"}}>
                                <Card>
                                    <CardTitle>Bugs</CardTitle>
                                    <CardBody>
                                        <Title headingLevel='h1' size="4xl">
                                            { apiDataCache[selected] &&
                                            <span>
                                                <span>{ apiDataCache[selected].resolution_time.resolved_bugs || "-"}</span>
                                            </span>
                                            }
                                            { !apiDataCache[selected] && "-" }
                                        </Title>
                                        <BugsChart chartType="bar" data={bugsChart}></BugsChart>
                                    </CardBody>
                                </Card>
                            </GridItem>
                        </Grid>
                    </GridItem>
                    <GridItem order={{default: "1"}}>
                        <Grid hasGutter span={3}>
                            <GridItem order={{default: "1"}}>
                                <Card isSelectable onClick={onClick} isSelected={selected.includes('Global')} id="Global">
                                    <CardTitle>Global</CardTitle>
                                    <CardBody>
                                        <Title headingLevel='h1' size="4xl">
                                        { apiDataCache["Global"] ? <span>{ apiDataCache["Global"].resolution_time.resolved_bugs + Math.floor(Math.random() * 205) } <span style={{fontSize: '10px'}}>open</span></span> : "-"}
                                        <br/>
                                        { apiDataCache["Global"] ? <span>{ apiDataCache["Global"].resolution_time.resolved_bugs} <span style={{fontSize: '10px'}}>resolved</span></span> : "-"}
                                        </Title>
                                    </CardBody>
                                </Card>
                            </GridItem>
                            <GridItem order={{default: "2"}}>
                                <Card isSelectable onClick={onClick} isSelected={selected.includes('Blocker')} id="Blocker">
                                    <CardTitle>Blockers</CardTitle>
                                    <CardBody>
                                        <Title headingLevel='h1' size="4xl">
                                        { apiDataCache["Blocker"] ? <span>{ apiDataCache["Blocker"].resolution_time.resolved_bugs + Math.floor(Math.random() * 205) } <span style={{fontSize: '10px'}}>open</span></span> : "-"}
                                        <br/>
                                        { apiDataCache["Blocker"] ? <span>{ apiDataCache["Blocker"].resolution_time.resolved_bugs} <span style={{fontSize: '10px'}}>resolved</span></span> : "-"}
                                        </Title>
                                    </CardBody>
                                </Card>
                            </GridItem>
                            <GridItem order={{default: "3"}}>
                                <Card isSelectable onClick={onClick} isSelected={selected.includes('Major')} id="Major">
                                    <CardTitle>Major Bugs</CardTitle>
                                    <CardBody>
                                        <Title headingLevel='h1' size="4xl">
                                        { apiDataCache["Major"] ? <span>{ apiDataCache["Major"].resolution_time.resolved_bugs + Math.floor(Math.random() * 205) } <span style={{fontSize: '10px'}}>open</span></span> : "-"}
                                        <br/>
                                        { apiDataCache["Major"] ? <span>{ apiDataCache["Major"].resolution_time.resolved_bugs} <span style={{fontSize: '10px'}}>resolved</span></span> : "-"}
                                        </Title>
                                    </CardBody>
                                </Card>
                            </GridItem>
                            <GridItem order={{default: "4"}}>
                                <Card isSelectable onClick={onClick} isSelected={selected.includes('Critical')} id="Critical">
                                    <CardTitle>Critical Bugs</CardTitle>
                                    <CardBody>
                                        <Title headingLevel='h1' size="4xl">
                                        { apiDataCache["Critical"] ? <span>{ apiDataCache["Critical"].resolution_time.resolved_bugs + Math.floor(Math.random() * 205) } <span style={{fontSize: '10px'}}>open</span></span> : "-"}
                                        <br/>
                                        { apiDataCache["Critical"] ? <span>{ apiDataCache["Critical"].resolution_time.resolved_bugs} <span style={{fontSize: '10px'}}>resolved</span></span> : "-"}
                                        </Title>
                                    </CardBody>
                                </Card>
                            </GridItem>
                            <GridItem order={{default: "5"}}>
                                <Card isSelectable onClick={onClick} isSelected={selected.includes('Normal')} id="Normal">
                                    <CardTitle>Normal Bugs</CardTitle>
                                    <CardBody>
                                        <Title headingLevel='h1' size="4xl">
                                        { apiDataCache["Normal"] ? <span>{ apiDataCache["Normal"].resolution_time.resolved_bugs + Math.floor(Math.random() * 205) } <span style={{fontSize: '10px'}}>open</span></span> : "-"}
                                        <br/>
                                        { apiDataCache["Normal"] ? <span>{ apiDataCache["Normal"].resolution_time.resolved_bugs} <span style={{fontSize: '10px'}}>resolved</span></span> : "-"}
                                        </Title>
                                    </CardBody>
                                </Card>
                            </GridItem>
                            <GridItem order={{default: "6"}}>
                                <Card isSelectable onClick={onClick} isSelected={selected.includes('Minor')} id="Minor">
                                    <CardTitle>Minor Bugs</CardTitle>
                                    <CardBody>
                                        <Title headingLevel='h1' size="4xl">
                                        { apiDataCache["Minor"] ? <span>{ apiDataCache["Minor"].resolution_time.resolved_bugs + Math.floor(Math.random() * 205) } <span style={{fontSize: '10px'}}>open</span></span> : "-"}
                                        <br/>
                                        { apiDataCache["Minor"] ? <span>{ apiDataCache["Minor"].resolution_time.resolved_bugs} <span style={{fontSize: '10px'}}>resolved</span></span> : "-"}
                                        </Title>
                                    </CardBody>
                                </Card>
                            </GridItem>
                            <GridItem order={{default: "6"}}>
                                <Card isSelectable onClick={onClick} isSelected={selected.includes('Undefined')} id="Undefined">
                                    <CardTitle>Undefined Bugs</CardTitle>
                                    <CardBody>
                                        <Title headingLevel='h1' size="4xl">
                                        { apiDataCache["Undefined"] ? <span>{ apiDataCache["Undefined"].resolution_time.resolved_bugs + Math.floor(Math.random() * 205) } <span style={{fontSize: '10px'}}>open</span></span> : "-"}
                                        <br/>
                                        { apiDataCache["Undefined"] ? <span>{ apiDataCache["Undefined"].resolution_time.resolved_bugs} <span style={{fontSize: '10px'}}>resolved</span></span> : "-"}
                                        </Title>
                                    </CardBody>
                                </Card>
                            </GridItem>
                        </Grid>
                    </GridItem>
                    <GridItem order={{default: "3"}}>
                        <Card style={{fontSize: "12px"}}>
                            <CardTitle>Bugs</CardTitle>
                            <CardBody>
                                <Grid hasGutter span={2}>
                                    <GridItem order={{default: "1"}}>
                                        <ToggleGroup aria-label="Default with single selectable">
                                            <ToggleGroupItem
                                                text="Resolved bugs"
                                                buttonId="resolved"
                                                isSelected={isSelected === 'resolved'}
                                                onChange={handleItemClick}
                                            />
                                            <ToggleGroupItem
                                                text="Opened bugs"
                                                buttonId="opened"
                                                isSelected={isSelected === 'opened'}
                                                onChange={handleItemClick}
                                            />
                                        </ToggleGroup>                                        
                                    </GridItem>
                                    <GridItem order={{default: "2"}}>
                                        <ChipGroup categoryName="Active filters: " numChips={5}>
                                            <Chip key={selected} isReadOnly style={{fontSize: '15px'}}>
                                                {selected} bugs
                                            </Chip>
                                        </ChipGroup>                                        
                                    </GridItem>
                                </Grid>
                                <ComposableTableStripedTr bugs={bugsTable}></ComposableTableStripedTr>
                            </CardBody>
                        </Card>
                    </GridItem>
                </Grid>
                </React.Fragment>
            </PageSection>
            
        </React.Fragment>
    )
}


const BugsChart: React.FC<{chartType:string, data:any}> = ({chartType, data}) => {

    const CursorVoronoiContainer = createContainer("voronoi", "cursor");

    let legendData: { name: string }[] = []
    if(data.length>0) {
        legendData = data.map((dataset, index) => { 
            return {name: dataset[0]["name"]}
        })
    } 
    
    return (
      <div style={{ margin: '0 auto', height: '60%', width: '90%', marginTop: '15px' }}>
        { data.length > 0 &&
        <Chart
          ariaDesc="Average number of pets"
          ariaTitle="Line chart example"
          height={210}
          legendData={legendData}
          legendPosition='bottom'          
          padding={{
            bottom: 70, // Adjusted to accommodate legend
            left:  40,
            right: 14,
            top: 20
          }}
        >
          <ChartAxis style={{ axisLabel: {fontSize: 8, padding: 30},tickLabels: {fontSize: 8}}}/>
          <ChartAxis dependentAxis={ true } showGrid style={{ axisLabel: {fontSize: 8, padding: 30}, tickLabels: {fontSize: 8}}}/>
          { chartType == 'bar' && data.length > 0 && 
          <ChartGroup offset={11}>
            {data.map((dataset, index) => (
                <ChartBar
                key={index}
                style={{
                    data: { strokeWidth: 1},
                    parent: { border: "1px solid #ccc"},
                    labels: { fill: "grey", fontSize: '7px' } 
                }}
                data={dataset}
                labels={({ datum }) => `${datum.y}`}
                />
            ))}
          </ChartGroup>
          }
          { chartType == 'line' && data.length > 0 && 
          <ChartGroup offset={11}>
            {data.map((dataset, index) => (
                <ChartLine
                key={index}
                style={{
                    data: { strokeWidth: 2},
                    parent: { border: "1px solid #ccc"},
                    labels: { fill: "grey", fontSize: '7px' } 
                }}
                data={dataset}
                labels={({ datum }) => `${parseInt(datum.y)}`}
                />
            ))}
          </ChartGroup>
          }
        </Chart>
        }
      </div>
    );
}

const ComposableTableStripedTr: React.FC<{bugs:any}> = ({bugs}) => {
    const [bugsPage, setBugsPage] = useState<Array<Bugs>>([]);
    const [page, setPage] = React.useState(1);
    const [perPage, setPerPage] = React.useState(10);

    useEffect(() => {

        if(bugs.length > 0){
            setBugsPage(bugs.slice(0, perPage))
            setPage(1)
        }
    }, [bugs]);

    const columnNames = {
        jira_key: "ID",
        created_at: "Created at",
        deleted_at: "Deleted at",
        updated_at: "Updated at",
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

    useEffect(() => {
        if(bugs.length > 0){
            let from = (page-1)*perPage
            let to = (page-1)*perPage + perPage > bugs.length ? bugs.length - 1 : (page-1)*perPage + perPage;
            setBugsPage(bugs.slice(from, to))
        }

    }, [page, perPage]);


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

  return (
    <div>
    <Pagination
        perPageComponent="button"
        itemCount={bugs.length}
        perPage={perPage}
        page={page}
        onSetPage={onSetPage}
        widgetId="top-example"
        onPerPageSelect={onPerPageSelect}
    />

    <TableComposable aria-label="Simple table" >
      <Caption>Jira bugs</Caption>
      <Thead>
        <Tr>
          <Th>{columnNames.jira_key}</Th>
          <Th>{columnNames.summary}</Th>
          <Th>{columnNames.status}</Th>
          <Th>{columnNames.created_at}</Th>
          <Th>{columnNames.updated_at}</Th>
        </Tr>
      </Thead>
      <Tbody>
        {bugsPage.map((bug, index) => (
          <Tr key={bug.jira_key} {...(index % 2 === 0 && { isStriped: true })}>
            <Td dataLabel={columnNames.jira_key}>{bug.jira_key}</Td>
            <Td dataLabel={columnNames.summary}>{bug.summary}</Td>
            <Td dataLabel={columnNames.status}>{bug.status}</Td>
            <Td dataLabel={columnNames.created_at}>{bug.created_at}</Td>
            <Td dataLabel={columnNames.updated_at}>{bug.updated_at}</Td>
          </Tr>
        ))}
      </Tbody>
    </TableComposable>

    <Pagination
        perPageComponent="button"
        itemCount={bugs.length}
        perPage={perPage}
        page={page}
        onSetPage={onSetPage}
        widgetId="top-example"
        onPerPageSelect={onPerPageSelect}
        />
    </div>
  );
};
