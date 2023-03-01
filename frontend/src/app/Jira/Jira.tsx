import React, { useContext, useEffect, useRef, useLayoutEffect, useState } from 'react';
import {
    Card,
    CardTitle,
    CardFooter,
    CardBody,
    Text,
    PageSection,
    Title,
    Grid,
    GridItem,
    TitleSizes,
    TextContent,
    PageSectionVariants,
    Pagination
} from '@patternfly/react-core';
import {
    TableComposable,
    Thead,
    Tr,
    Th,
    Tbody,
    Td
} from '@patternfly/react-table';

import { Chart, ChartAxis, ChartGroup, ChartLine,ChartBar, ChartLegendTooltip, createContainer } from '@patternfly/react-charts';

import { Caption} from '@patternfly/react-table';
import { Chip, ChipGroup } from '@patternfly/react-core';

import { ChartDonut, ChartThemeColor } from '@patternfly/react-charts';
import { getJiras } from '@app/utils/APIService';
import { ReactReduxContext } from 'react-redux';
import { isValidTeam } from '@app/utils/utils';

export const Jira = () => {

    const [selected, setSelected] = useState<Array<string>>([]);

    const onClick = (event: React.MouseEvent) => {
        if(selected.includes(event.currentTarget.id)){
            var array = [...selected];
            const index = selected.indexOf(event.currentTarget.id);
            if (index !== -1) {
                array.splice(index, 1);
                setSelected(array)
            }
        } else {
            setSelected([...selected, event.currentTarget.id])
        }

    };

    const deleteItem = (id: string) => {
        if(selected.includes(id)){
            var array = [...selected];
            const index = selected.indexOf(id);
            if (index !== -1) {
                array.splice(index, 1);
                setSelected(array)
            }
        }
    };
    
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
                        <Grid hasGutter sm={6} md={6} lg={6} xl={6} xl2={6}>
                            <GridItem order={{default: "1"}}>
                                <Card>
                                    <CardTitle>Average Resolution Time</CardTitle>
                                    <CardBody>
                                        <Title headingLevel='h1' size="4xl">3.64 years</Title>
                                        <BugsChart chartType="line"></BugsChart>
                                    </CardBody>
                                </Card>
                            </GridItem>
                            <GridItem order={{default: "2"}}>
                                <Card>
                                    <CardTitle>Open Bugs</CardTitle>
                                    <CardBody>
                                        <Title headingLevel='h1' size="4xl">100k</Title>
                                        <BugsChart chartType="bar"></BugsChart>
                                    </CardBody>
                                </Card>
                            </GridItem>
                        </Grid>
                    </GridItem>
                    <GridItem order={{default: "1"}}>
                        <Grid hasGutter span={3}>
                            <GridItem order={{default: "3"}}>
                                <Card isSelectable onClick={onClick} isSelected={selected.includes('jiras')} id="jiras">
                                    <CardTitle>Jiras</CardTitle>
                                    <CardBody><Title headingLevel='h1' size="4xl">0</Title></CardBody>
                                </Card>
                            </GridItem>
                            <GridItem order={{default: "4"}}>
                                <Card isSelectable onClick={onClick} isSelected={selected.includes('blockers')} id="blockers">
                                    <CardTitle>Blockers</CardTitle>
                                    <CardBody><Title headingLevel='h1' size="4xl">23</Title></CardBody>
                                </Card>
                            </GridItem>
                            <GridItem order={{default: "5"}}>
                                <Card isSelectable onClick={onClick} isSelected={selected.includes('major')} id="major">
                                    <CardTitle>Major Bugs</CardTitle>
                                    <CardBody><Title headingLevel='h1' size="4xl">1</Title></CardBody>
                                </Card>
                            </GridItem>
                            <GridItem order={{default: "6"}}>
                                <Card isSelectable onClick={onClick} isSelected={selected.includes('critical')} id="critical">
                                    <CardTitle>Critical Bugs</CardTitle>
                                    <CardBody><Title headingLevel='h1' size="4xl">89</Title></CardBody>
                                </Card>
                            </GridItem>
                        </Grid>
                    </GridItem>
                    <GridItem order={{default: "3"}}>
                        <Card style={{fontSize: "12px"}}>
                            <CardTitle>Bugs</CardTitle>
                            <CardBody>
                            <ChipGroup categoryName="Bugs filters: " numChips={5}>
                                {selected.map(currentChip => (
                                    <Chip key={currentChip} onClick={() => deleteItem(currentChip)} style={{fontSize: '15px'}}>
                                    {currentChip}
                                    </Chip>
                                ))}
                            </ChipGroup>
                            <ComposableTableStripedTr ></ComposableTableStripedTr>
                            </CardBody>
                        </Card>
                    </GridItem>
                </Grid>
                </React.Fragment>
            </PageSection>
        </React.Fragment>
    )
}


const BugsChart: React.FC<{chartType:string}> = ({chartType}) => {

    // Note: Container order is important
    const CursorVoronoiContainer = createContainer("voronoi", "cursor");

    return (
      <div style={{ margin: '0 auto', height: '60%', width: '90%', marginTop: '15px' }}>
        <Chart
          ariaDesc="Average number of pets"
          ariaTitle="Line chart example"
          maxDomain={{y: 100}}
          minDomain={{y: 0}}
          height={210}
          name="chart2"
          padding={{
            bottom: 40, // Adjusted to accommodate legend
            left:  14,
            right: 14,
            top: 10
          }}
        >
          <ChartAxis style={{ axisLabel: {fontSize: 8, padding: 30},tickLabels: {fontSize: 8}}}/>
          <ChartAxis showGrid style={{ axisLabel: {fontSize: 8, padding: 30}, tickLabels: {fontSize: 8}}}/>
          { chartType == 'bar' && <ChartGroup offset={11}>
            <ChartBar
              style={{
                data: { strokeWidth: 1},
                parent: { border: "1px solid #ccc"}
              }}
              data={[
                { name: "Open bugs", x: 'Jan', y: 1 },
                { name: "Open bugs", x: 'Feb', y: 2 },
                { name: "Open bugs", x: 'Mar', y: 5 },
                { name: "Open bugs", x: 'Apr', y: 3 },
                { name: "Open bugs", x: 'May', y: 3 },
                { name: "Open bugs", x: 'Jun', y: 3 }
              ]}
            />
            <ChartBar
              style={{
                data: { strokeWidth: 1},
                parent: { border: "1px solid #ccc"}
              }}
              data={[
                { name: "Total bugs", x: 'Jan', y: 0 },
                { name: "Total bugs", x: 'Feb', y: 6 },
                { name: "Total bugs", x: 'Mar', y: 89 },
                { name: "Total bugs", x: 'Apr', y: 5 },
                { name: "Total bugs", x: 'May', y: 45 }
              ]}
            />
          </ChartGroup>
          }
          { chartType == 'line' && <ChartGroup offset={11}>
            <ChartLine
              style={{
                data: { strokeWidth: 1},
                parent: { border: "1px solid #ccc"}
              }}
              data={[
                { name: "Open bugs", x: 'Jan', y: 1 },
                { name: "Open bugs", x: 'Feb', y: 2 },
                { name: "Open bugs", x: 'Mar', y: 5 },
                { name: "Open bugs", x: 'Apr', y: 3 },
                { name: "Open bugs", x: 'May', y: 3 },
                { name: "Open bugs", x: 'Jun', y: 3 },
                { name: "Open bugs", x: 'Jul', y: 3 },
                { name: "Open bugs", x: 'Aug', y: 3 },
                { name: "Open bugs", x: 'Sep', y: 3 },
                { name: "Open bugs", x: 'Oct', y: 3 },
                { name: "Open bugs", x: 'Nov', y: 3 },
                { name: "Open bugs", x: 'Dec', y: 3 }
              ]}
            />
            <ChartLine
              style={{
                data: { strokeWidth: 1},
                parent: { border: "1px solid #ccc"}
              }}
              data={[
                { name: "Total bugs", x: 'Jan', y: 0 },
                { name: "Total bugs", x: 'Feb', y: 6 },
                { name: "Total bugs", x: 'Mar', y: 89 },
                { name: "Total bugs", x: 'Apr', y: 5 },
                { name: "Total bugs", x: 'May', y: 45 },
                { name: "Total bugs", x: 'Jun', y: 40 },
                { name: "Total bugs", x: 'Jul', y: 30 },
                { name: "Total bugs", x: 'Aug', y: 34 },
                { name: "Total bugs", x: 'Sep', y: 29 },
                { name: "Total bugs", x: 'Oct', y: 6 },
                { name: "Total bugs", x: 'Nov', y: 9 },
                { name: "Total bugs", x: 'Dec', y: 0 }
              ]}
            />
          </ChartGroup>
          }
        </Chart>
      </div>
    );
}


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

export const ComposableTableStripedTr: React.FunctionComponent = () => {
    const { store } = useContext(ReactReduxContext);
    const state = store.getState();
    const dispatch = store.dispatch;
    const [bugs, setBugs] = useState<Array<Bugs>>([]);
    
    useEffect(() => {
        getJiras().then((res) => {
            if (res.code === 200) {
                const result = res.data;
                dispatch({ type: "SET_JIRAS", data: result });
                setBugs(result.slice(0, perPage))
            } else {
                dispatch({ type: "SET_ERROR", data: res });
            }
        })
    }, []);

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

  const [page, setPage] = React.useState(1);
  const [perPage, setPerPage] = React.useState(5);

  useEffect(() => {
    let from = (page-1)*perPage
    let to = (page-1)*perPage + perPage > state.jiras["E2E_KNOWN_ISSUES"].length ? state.jiras["E2E_KNOWN_ISSUES"].length - 1 : (page-1)*perPage + perPage;
    console.log(from, to, state.jiras["E2E_KNOWN_ISSUES"].length - 1)
    setBugs(state.jiras["E2E_KNOWN_ISSUES"].slice(from, to))
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
        itemCount={state.jiras["E2E_KNOWN_ISSUES"].length}
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
        {bugs.map((bug, index) => (
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
        itemCount={state.jiras["E2E_KNOWN_ISSUES"].length}
        perPage={perPage}
        page={page}
        onSetPage={onSetPage}
        widgetId="top-example"
        onPerPageSelect={onPerPageSelect}
        />
    </div>
  );
};
