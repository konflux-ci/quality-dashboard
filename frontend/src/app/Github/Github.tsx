import React, { useEffect, useState } from 'react';
import { CopyIcon, ExclamationCircleIcon, PlusIcon } from '@patternfly/react-icons';
import {
    PageSection,
    PageSectionVariants,
    Title,
    TitleSizes,
    ButtonVariant,
    EmptyState,
    EmptyStateIcon,
    EmptyStateVariant,
} from '@patternfly/react-core';
import { Toolbar, ToolbarItem, ToolbarContent } from '@patternfly/react-core';
import { Button } from '@patternfly/react-core';
import { Select, SelectOption, SelectVariant } from '@patternfly/react-core';
import { deleteInApi, getAllRepositoriesWithOrgs, getPullRequests, getWorkflowByRepositoryName } from '@app/utils/APIService';
import { Grid, GridItem } from '@patternfly/react-core';
import {
    InfoCard,
} from '@app/utils/sharedComponents';
import { ReactReduxContext, useSelector } from 'react-redux';
import { useHistory } from 'react-router-dom';
import { CodeCov } from './CodeCov';
import { GithubActions } from './GithubActions';
import { ActionsColumn, IAction } from '@patternfly/react-table';
import { FormModal, ModalContext, useDefaultModalContextState, useModalContext } from './CreateRepository';
import { isValidTeam } from '@app/utils/utils';
import { PrsStatistics, PullRequestCard, PullRequestsGraphic } from './PullRequests';
import { formatDate, getRangeDates } from '@app/Reports/utils';
import { DateTimeRangePicker } from '@app/utils/DateTimeRangePicker';


// eslint-disable-next-line prefer-const
let GitHub = () => {

    const { store } = React.useContext(ReactReduxContext);
    const state = store.getState();
    const [repositories, setRepositories] = useState<{ repoName: string, organization: string, description: string, codeCoverage: any, isPlaceholder?: boolean }[]>([]);
    const [repoName, setRepoName] = useState("");
    const [repoOrg, setRepoOrg] = useState("");
    const [repoNameToggle, setRepoNameToggle] = useState(false);
    const [workflows, setWorkflows] = useState([]);
    const [prs, setPrs] = useState<PrsStatistics | null>(null);
    const [rangeDateTime, setRangeDateTime] = useState(getRangeDates(365));
    const [noData, setNoData] = useState(false)
    const defaultModalContext = useDefaultModalContextState();
    const modalContext = useModalContext()
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

    // Reset rangeDateTime
    const clearRangeDateTime = () => {
        setRangeDateTime(getRangeDates(365))
    }

    // Reset all dropwdowns and state variables
    const clearAll = () => {
        clearRepo()
        clearRangeDateTime()
        setNoData(false)
    }

    // Reset params
    const clearParams = () => {
        clearAll()
        history.push(window.location.pathname + '?' + "team=" + params.get("team"));
    }

    const validatePullRequests = () => {
        setNoData(false)
        if (repositories.find(r => r.organization == repoOrg && r.repoName == repoName)) {
            try {
                getPullRequests(repoName, repoOrg, rangeDateTime)
                    .then((data: any) => {
                        setPrs(data)
                    });
            } catch(error) {
                console.log(error)
                setNoData(true)
            }
        } else {
            setNoData(true)
        }
    }

    // Triggers automatic validation when state variables change
    useEffect(() => {
        validatePullRequests();
    }, [repoOrg, repoName, rangeDateTime]);


    useEffect(() => {
        if (state.teams.Team != "") {
            setRepositories([])
            clearAll()

            const repository = params.get("repository")
            const organization = params.get("organization")
            const team = params.get("team")
            const start = params.get("start")
            const end = params.get("end")

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

                        if (repository == null || organization == null || start == null || end == null) { // first click on page or team
                            setRepoName(data[1].repoName)
                            setRepoOrg(data[1].organization)
                            getWorkflows(data[1].repoName)

                            const start_date = formatDate(rangeDateTime[0])
                            const end_date = formatDate(rangeDateTime[1])

                            history.push('/home/github?team=' + currentTeam + '&organization=' + data[1].organization + '&repository=' + data[1].repoName
                                + '&start=' + start_date + '&end=' + end_date)

                        } else {
                            setRepoName(repository)
                            setRepoOrg(organization)
                            getWorkflows(repository)
                            setRangeDateTime([new Date(start), new Date(end)])

                            history.push('/home/github?team=' + currentTeam + '&organization=' + organization + '&repository=' + repository
                                + '&start=' + start + '&end=' + end)
                        }
                    }
                })
        }
    }, [setRepositories, currentTeam]);


    async function deleteRepository(gitOrg: string, repoName: string) {
        const data = {
            git_organization: gitOrg,
            repository_name: repoName,
        }
        try {
            await deleteInApi(data, '/api/quality/repositories/delete')
            history.push(window.location.pathname + '?' + "team=" + params.get("team"));
            window.location.reload();
        } catch (error) {
            console.log(error)
        }
    }


    async function editRepository(repo) {
        try {
            modalContext.handleModalToggle(true, repo)
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


    // Validates if the repository and  organization are correct
    const validQueryParams = (repository, organization) => {
        if (isValidTeam()) {
            if (repositories.find(r => r.organization == organization && r.repoName == repository)) {
                return true;
            }
            if (repository == "" && organization == "") {
                return true;
            }
        }
        return false;
    }

    function handleChange(event, from, to) {
        setRangeDateTime([from, to])
        params.set("start", formatDate(from))
        params.set("end", formatDate(to))
        history.push(window.location.pathname + '?' + params.toString());
    }

    const start = rangeDateTime[0]
    const end = rangeDateTime[1]

    return (
        <ModalContext.Provider value={defaultModalContext}>
            <React.Fragment>
                {/* page title bar */}
                <PageSection variant={PageSectionVariants.light}>
                    <Title headingLevel="h3" size={TitleSizes['2xl']}>
                        GitHub metrics
                        <Button onClick={() => navigator.clipboard.writeText(window.location.href)} variant="link" icon={<CopyIcon />} iconPosition="right">
                            Copy link
                        </Button>
                        <Button style={{ float: 'right' }} variant={ButtonVariant.secondary} onClick={modalContext.handleModalToggle}>
                            <PlusIcon /> &nbsp; Add a repository
                        </Button>
                    </Title>
                </PageSection>
                {/* main content  */}
                <PageSection>
                    {/* the following toolbar will contain the form (dropdowns and button) to request data to the server */}
                    <Toolbar style={{ width: '100%', margin: 'auto' }}>
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
                            <ActionsColumn items={defaultActions(getRepository(repoName, repoOrg))} />
                            <ToolbarItem style={{ minWidth: "fitContent", maxWidth: "fitContent" }}>
                                <DateTimeRangePicker
                                    startDate={start}
                                    endDate={end}
                                    handleChange={(event, from, to) => handleChange(event, from, to)}
                                >
                                </DateTimeRangePicker>
                            </ToolbarItem>
                            <ToolbarItem style={{ minWidth: "fitContent", maxWidth: "fitContent" }}>
                                <Button variant="link" onClick={clearParams}>Clear</Button>
                            </ToolbarItem>
                            <FormModal></FormModal>
                        </ToolbarContent>
                    </Toolbar>

                    {/* this section will show statistics and details about GitHub metric */}
                    {validQueryParams(repoName, repoOrg) && prs != null && prs.metrics != undefined && !noData && <div style={{ marginTop: '20px' }}>
                        <Grid hasGutter>
                            <GridItem span={3} rowSpan={2}>
                                <InfoCard data={[{ title: "Repository", value: repoName }, { title: "Organization", value: repoOrg }, { title: "Description", value: getDescription(repoName, repoOrg) }]}></InfoCard>
                            </GridItem>

                            <GridItem span={7} rowSpan={4}>
                                <PullRequestsGraphic metrics={prs?.metrics}></PullRequestsGraphic>
                            </GridItem>

                            <GridItem span={2} rowSpan={1}>
                                <PullRequestCard title="Open PRs" total={prs?.summary?.open_prs}></PullRequestCard>
                            </GridItem>

                            <GridItem span={2} rowSpan={1}>
                                <PullRequestCard title="Merged PRs" total={prs?.summary?.merged_prs}></PullRequestCard>
                            </GridItem>

                            <GridItem span={3} rowSpan={2}>
                                <CodeCov repo={getRepository(repoName, repoOrg)}></CodeCov>
                            </GridItem>

                            <GridItem span={2} rowSpan={2}>
                                <PullRequestCard title="Merge PR Avg Days" total={prs?.summary?.merge_avg}></PullRequestCard>
                            </GridItem>
{/*                             
                            <GridItem span={2} rowSpan={1}>
                                <PullRequestCard title="Merge PR Avg Days" total={prs?.summary.merge_avg}></PullRequestCard>
                            </GridItem> */}

                            {workflows.length > 0 && <GridItem span={12}>
                                <GithubActions repoName={repoName} workflows={workflows}></GithubActions>
                            </GridItem>
                            }
                        </Grid>
                    </div >
                    }
                    {!validQueryParams(repoName, repoOrg) && <EmptyState variant={EmptyStateVariant.xl}>
                        <EmptyStateIcon icon={ExclamationCircleIcon} />
                        <Title headingLevel="h1" size="lg">
                            Something went wrong.
                        </Title>
                    </EmptyState>
                    }
                    {noData && <EmptyState variant={EmptyStateVariant.xl}>
                        <EmptyStateIcon icon={ExclamationCircleIcon} />
                        <Title headingLevel="h1" size="lg">
                            No repository detected.
                        </Title>
                    </EmptyState>
                    }
                </PageSection >
            </React.Fragment >
        </ModalContext.Provider>
    )
}

export { GitHub };

