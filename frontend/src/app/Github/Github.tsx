import React, { useEffect, useState } from 'react';
import { CopyIcon, PlusIcon } from '@patternfly/react-icons';
import {
    PageSection,
    PageSectionVariants,
    Title,
    TitleSizes,
    CardBody,
    Card,
    ButtonVariant,
} from '@patternfly/react-core';
import { Toolbar, ToolbarItem, ToolbarContent } from '@patternfly/react-core';
import { Button } from '@patternfly/react-core';
import { Select, SelectOption, SelectVariant } from '@patternfly/react-core';
import { deleteInApi, getAllRepositoriesWithOrgs, getWorkflowByRepositoryName } from '@app/utils/APIService';
import { Grid, GridItem } from '@patternfly/react-core';
import {
    InfoCard,
} from '@app/utils/sharedComponents';
import { ReactReduxContext, useSelector } from 'react-redux';
import { useHistory } from 'react-router-dom';
import { Chart, ChartAxis, ChartBar, ChartGroup, ChartVoronoiContainer } from '@patternfly/react-charts';
import { CodeCov } from './CodeCov';
import { GithubActions } from './GithubActions';
import { ActionsColumn, IAction } from '@patternfly/react-table';


// eslint-disable-next-line prefer-const
let GitHub = () => {

    const { store } = React.useContext(ReactReduxContext);
    const state = store.getState();
    const [repositories, setRepositories] = useState<{ repoName: string, organization: string, description: string, codeCoverage: any, isPlaceholder?: boolean }[]>([]);
    const [repoName, setRepoName] = useState("");
    const [repoOrg, setRepoOrg] = useState("");
    const [repoNameToggle, setRepoNameToggle] = useState(false);
    const [workflows, setWorkflows] = useState([]);
    const [isOpen, setOpen] = useState<boolean>(false);
    const currentTeam = useSelector((state: any) => state.teams.Team);
    const history = useHistory();
    const params = new URLSearchParams(window.location.search);

    function getWorkflows(repo) {
        getWorkflowByRepositoryName(repo).then((res) => {
            if (res.code === 200) {
                const result = res.data.sort((a, b) => (a.workflow_name < b.workflow_name ? -1 : 1));
                setWorkflows(result)
            }
        });
    }

    const getDescription = (repoName, repoOrg) => {
        const repo = getRepository(repoName, repoOrg)

        if (repo != undefined) {
            return repo.description
        }

        return ""
    }

    const getRepository = (repoName, repoOrg) => {
        const repo = repositories.find(r => r.organization == repoOrg && r.repoName == repoName)

        return repo
    }

    // Reset the repository dropdown
    const clearRepo = () => {
        setRepoName("")
        setRepoOrg("")
        setRepoNameToggle(false)
        setWorkflows([])
    }

    // Called onChange of the repository dropdown element. This set repository name and organization state variables, or clears them when placeholder is selected
    const setRepoNameOnChange = (event, selection, isPlaceholder) => {
        if (isPlaceholder) {
            clearRepo()
        }
        else {
            setRepoName(repositories[selection].repoName);
            setRepoOrg(repositories[selection].organization);
            setRepoNameToggle(false)
            getWorkflows(repositories[selection].repoName)
            params.set("repository", repositories[selection].repoName)
            params.set("organization", repositories[selection].organization)
            history.push(window.location.pathname + '?' + params.toString());
        };
    }

    // Reset all dropwdowns and state variables
    const clearAll = () => {
        clearRepo()
    }

    // Reset params
    const clearParams = () => {
        clearAll()
        history.push(window.location.pathname + '?' + "team=" + params.get("team"));
    }


    useEffect(() => {
        if (state.teams.Team != "") {
            setRepositories([])
            clearAll()

            const repository = params.get("repository")
            const organization = params.get("organization")
            const team = params.get("team")

            getAllRepositoriesWithOrgs(state.teams.Team, false)
                .then((data: any) => {
                    let dropDescr = ""
                    if (data.length < 1 && (team == state.teams.Team || team == null)) {
                        dropDescr = "No Repositories"
                        history.push('/home/github?team=' + currentTeam)
                    }
                    else { dropDescr = "Select a repository" }

                    if (data.length > 0 && (team == state.teams.Team || team == null)) {
                        data.unshift({ repoName: dropDescr, organization: "", isPlaceholder: true }) // Adds placeholder at the beginning of the array, so it will be shown first
                        setRepositories(data)

                        if (repository == null || organization == null) { // first click on page or team
                            setRepoName(data[1].repoName)
                            setRepoOrg(data[1].organization)
                            getWorkflows(data[1].repoName)
                            history.push('/home/github?team=' + currentTeam + '&organization=' + data[1].organization + '&repository=' + data[1].repoName)

                        } else {
                            setRepoName(repository)
                            setRepoOrg(organization)
                            getWorkflows(repository)

                            history.push('/home/github?team=' + currentTeam + '&organization=' + organization + '&repository=' + repository)
                        }
                    }
                })
        }
    }, [setRepositories, currentTeam]);


    const handleModalToggle = () => {
        setOpen(!isOpen)
    };

    async function deleteRepository(gitOrg: string, repoName: string) {
        const data = {
            git_organization: gitOrg,
            repository_name: repoName,
        }
        try {
            await deleteInApi(data, '/api/quality/repositories/delete')
        } catch (error) {
            console.log(error)
        }
        history.push(window.location.pathname + '?' + "team=" + params.get("team"));
        window.location.reload();
    }


    async function editRepository(repo) {
        try {
            console.log("ioi")
        } catch (error) {
            console.log(error)
        }
    }


    const defaultActions = (repo): IAction[] => [
        {
            title: 'Delete Repository',
            onClick: () => deleteRepository(repo.organization, repo.repoName)
        },
        {
            title: 'Edit Repository',
            onClick: () => editRepository(repo)
        },
    ];


    return (

        <React.Fragment>
            {/* page title bar */}
            <PageSection variant={PageSectionVariants.light}>
                <Title headingLevel="h3" size={TitleSizes['2xl']}>
                    GitHub metrics
                    <Button onClick={() => navigator.clipboard.writeText(window.location.href)} variant="link" icon={<CopyIcon />} iconPosition="right">
                        Copy link
                    </Button>
                    <Button style={{ float: 'right' }} variant={ButtonVariant.secondary} onClick={handleModalToggle}>
                            <PlusIcon /> &nbsp; Add a repository
                    </Button>
                </Title>
            </PageSection>
            {/* main content  */}
            <PageSection>
                {/* the following toolbar will contain the form (dropdowns and button) to request data to the server */}
                <Toolbar style={{ width: '100%', margin: 'auto'}}>
                    <ToolbarContent style={{ textAlign: 'center' }}>
                        <ToolbarItem style={{ minWidth: "20%", maxWidth: "40%" }}>
                            <span id="typeahead-select" hidden>
                                Select a state
                            </span>
                            <Select
                                variant={SelectVariant.typeahead}
                                typeAheadAriaLabel="Select a repository"
                                isOpen={repoNameToggle}
                                onToggle={setRepoNameToggle}
                                selections={repoName}
                                onSelect={setRepoNameOnChange}
                                onClear={clearRepo}
                                aria-labelledby="typeahead-select"
                                placeholderText="Select a repository"
                            >
                                {repositories.map((value, index) => (
                                    <SelectOption key={index} value={index} description={value.organization} isDisabled={value.isPlaceholder}>{value.repoName}</SelectOption>
                                ))}
                            </Select>
                        </ToolbarItem>
                        <ActionsColumn items={defaultActions(getRepository(repoName, repoOrg))}/>
                        <ToolbarItem style={{ minWidth: "fitContent", maxWidth: "fitContent" }}>
                            <Button variant="link" onClick={clearParams}>Clear</Button>
                        </ToolbarItem>
                    </ToolbarContent>
                </Toolbar>
        
                {/* this section will show statistics and details about GitHub metric */}
                <div style={{ marginTop: '20px' }}>
                    <Grid hasGutter>
                        <GridItem span={3} rowSpan={2}>
                            <InfoCard data={[{ title: "Repository", value: repoName }, { title: "Organization", value: repoOrg }, { title: "Description", value: getDescription(repoName, repoOrg) }]}></InfoCard>
                        </GridItem>


                        <GridItem span={7} rowSpan={4}>
                            <Card>
                                <CardBody>
                                    <Title headingLevel={'h2'}>Pull Requests over time</Title>
                                    <Chart
                                        ariaDesc="Average number of pets"
                                        ariaTitle="Bar chart example"
                                        containerComponent={<ChartVoronoiContainer labels={({ datum }) => `${datum.name}: ${datum.y}`} constrainToVisibleArea />}
                                        domain={{ y: [0, 9] }}
                                        domainPadding={{ x:  [30, 25] }}
                                        legendData={[{ name: 'Cats' }, { name: 'Dogs' }, { name: 'Birds' }, { name: 'Mice' }]}
                                        legendOrientation="vertical"
                                        legendPosition="right"
                                        height={250}
                                        name="chart1"
                                        padding={{
                                            bottom: 50,
                                            left: 50,
                                            right: 200, // Adjusted to accommodate legend
                                            top: 50
                                        }}
                                        width={600}
                                    >
                                        <ChartAxis />
                                        <ChartAxis dependentAxis showGrid />
                                        <ChartGroup offset={11}>
                                            <ChartBar data={[{ name: 'Cats', x: '2015', y: 1 }, { name: 'Cats', x: '2016', y: 2 }, { name: 'Cats', x: '2017', y: 5 }, { name: 'Cats', x: '2018', y: 3 }]} />
                                            <ChartBar data={[{ name: 'Dogs', x: '2015', y: 2 }, { name: 'Dogs', x: '2016', y: 1 }, { name: 'Dogs', x: '2017', y: 7 }, { name: 'Dogs', x: '2018', y: 4 }]} />
                                            <ChartBar data={[{ name: 'Birds', x: '2015', y: 4 }, { name: 'Birds', x: '2016', y: 4 }, { name: 'Birds', x: '2017', y: 9 }, { name: 'Birds', x: '2018', y: 7 }]} />
                                            <ChartBar data={[{ name: 'Mice', x: '2015', y: 3 }, { name: 'Mice', x: '2016', y: 3 }, { name: 'Mice', x: '2017', y: 8 }, { name: 'Mice', x: '2018', y: 5 }]} />
                                        </ChartGroup>
                                    </Chart>
                                </CardBody>
                            </Card>
                        </GridItem>

                        <GridItem span={2} rowSpan={2}>
                            <InfoCard data={[{ title: "Open Prs", value: repoName }]}></InfoCard>
                        </GridItem>

                        <GridItem span={3} rowSpan={2}>
                            <CodeCov repo={getRepository(repoName, repoOrg)}></CodeCov>
                        </GridItem>

                        <GridItem span={2} rowSpan={2}>
                            <InfoCard data={[{ title: "Closed Prs", value: repoName }]}></InfoCard>
                        </GridItem>

                        {workflows.length > 0 && <GridItem span={12}>
                            <GithubActions repoName={repoName} workflows={workflows}></GithubActions>
                        </GridItem>
                        }
                    </Grid>
                </div >
            </PageSection >
        </React.Fragment >
    )
}

export { GitHub };

