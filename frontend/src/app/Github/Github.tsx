import React, { useEffect, useState } from 'react';
import { CopyIcon, ExclamationCircleIcon, PlusIcon } from '@patternfly/react-icons';
import {
  PageSection,
  PageSectionVariants,
  Title,
  TitleSizes,
  ButtonVariant,
  Spinner,
  Card,
  CardTitle,
  CardBody,
  EmptyState,
  EmptyStateIcon,
  EmptyStateVariant,
  Flex,
  FlexItem,
} from '@patternfly/react-core';
import { Button } from '@patternfly/react-core';
import {
  getAllRepositoriesWithOrgs,
  getTeams,
} from '@app/utils/APIService';
import { Grid, GridItem } from '@patternfly/react-core';
import { ReactReduxContext, useSelector } from 'react-redux';
import { useHistory } from 'react-router-dom';
import { Workflows } from './GithubActions';
import { FormModal, ModalContext, useDefaultModalContextState, useModalContext } from './CreateRepository';
import { validateParam } from '@app/utils/utils';
import { PrsStatistics } from './PullRequests';
import { formatDate, getRangeDates } from '@app/Reports/utils';
import { DateTimeRangePicker } from '@app/utils/DateTimeRangePicker';
import { ComposableTable } from './Table';
import { Coverage } from './CodeCov';
import { Header } from '@app/utils/Header';

export interface RepositoryInfo {
  repository_name: string;
  git_organization: string;
  description: string;
  git_url: string;
  code_coverage: Coverage;
  prs: PrsStatistics;
  workflows: Workflows[];
  code_cov: string;
  coverage_trend: string,
  retest_avg: string;
  retest_before_merge_avg: string;
  created_prs_in_time_range: string;
  open_prs: string;
  merged_prs: string;
  merged_prs_in_time_range: string;
  time_to_merge_pr_avg_days: string;
}

// eslint-disable-next-line prefer-const
let GitHub = () => {
  const { store } = React.useContext(ReactReduxContext);
  const state = store.getState();
  const [rangeDateTime, setRangeDateTime] = useState(getRangeDates(90));
  const defaultModalContext = useDefaultModalContextState();
  const modalContext = useModalContext();
  const currentTeam = useSelector((state: any) => state.teams.Team);
  const history = useHistory();
  const params = new URLSearchParams(window.location.search);
  const [loadingState, setLoadingState] = useState(false);
  const [repos, setRepos] = useState<any>({});
  const [isInvalid, setIsInvalid] = useState(false);


  // Reset rangeDateTime
  const clearRangeDateTime = () => {
    setRangeDateTime(getRangeDates(90));
  };

  // Triggers automatic validation when state variables change
  useEffect(() => {
    getAllRepositoriesWithOrgs(state.teams.Team, false, rangeDateTime).then((data: any) => {
      const rps = new Array<RepositoryInfo>
      data.forEach((repository, _) => {
        rps.push({
          repository_name: repository.repoName,
          git_organization: repository.organization,
          description: repository.description,
          git_url: repository.url,
          code_coverage: repository.coverage,
          prs: repository.prs,
          workflows: repository.workflows,
          code_cov: repository.coverage.coverage_percentage == 0 ? 'N/A' : repository.coverage.coverage_percentage,
          coverage_trend: repository.coverage.coverage_trend,
          // edge case of service-provider-integration-operator
          // https://github.com/redhat-appstudio/service-provider-integration-operator/pull/548#issuecomment-1494149514
          retest_avg: repository.prs?.summary.retest_avg == 0 || repository.prs?.summary.retest_avg == 0.01 ? 'N/A' : repository.prs?.summary.retest_avg,
          retest_before_merge_avg: repository.prs?.summary.retest_before_merge_avg == 0 || repository.prs?.summary.retest_before_merge_avg == 0.01 ? 'N/A' : repository.prs?.summary.retest_before_merge_avg,
          created_prs_in_time_range: repository.prs?.summary?.created_prs_in_time_range,
          open_prs: repository.prs?.summary?.open_prs,
          merged_prs: repository.prs?.summary?.merged_prs,
          merged_prs_in_time_range: repository.prs?.summary?.merged_prs_in_time_range,
          time_to_merge_pr_avg_days: repository.prs?.summary?.merge_avg + ' day(s)',
        });
      })
      setRepos(rps)
    });
  }, [rangeDateTime]);


  useEffect(() => {
    setLoadingState(true)
    setIsInvalid(false);

    const team = params.get("team")
    if ((team != null) && (team != state.teams.Team)) {
      getTeams().then(res => {
        if (!validateParam(res.data, team)) {
          setLoadingState(false)
          setIsInvalid(true)
        }
      })
    }

    if (state.teams.Team != '') {
      setRepos([]);

      const team = params.get('team');
      const start = params.get('start');
      const end = params.get('end');

      getAllRepositoriesWithOrgs(state.teams.Team, false, rangeDateTime).then((data: any) => {
        if (data.length < 1 && (team == state.teams.Team || team == null)) {
          setLoadingState(false)
          history.push('/home/github?team=' + currentTeam);
        }

        if (data.length > 0 && (team == state.teams.Team || team == null)) {
          let rps = new Array<RepositoryInfo>
          data.forEach((repository, _) => {
            rps.push({
              repository_name: repository.repoName,
              git_organization: repository.organization,
              description: repository.description,
              git_url: repository.url,
              code_coverage: repository.coverage,
              prs: repository.prs,
              workflows: repository.workflows,
              code_cov: repository.coverage.coverage_percentage == 0 ? 'N/A' : repository.coverage.coverage_percentage,
              coverage_trend: repository.coverage.coverage_trend,
              // edge case of service-provider-integration-operator
              // https://github.com/redhat-appstudio/service-provider-integration-operator/pull/548#issuecomment-1494149514
              retest_avg: repository.prs?.summary.retes_avg == 0 || repository.prs?.summary.retest_avg == 0.01 ? 'N/A' : repository.prs?.summary.retest_avg,
              retest_before_merge_avg: repository.prs?.summary.retest_before_merge_avg == 0 || repository.prs?.summary.retest_before_merge_avg == 0.01 ? 'N/A' : repository.prs?.summary.retest_before_merge_avg,
              created_prs_in_time_range: repository.prs?.summary?.created_prs_in_time_range,
              open_prs: repository.prs?.summary?.open_prs,
              merged_prs: repository.prs?.summary?.merged_prs,
              merged_prs_in_time_range: repository.prs?.summary?.merged_prs_in_time_range,
              time_to_merge_pr_avg_days: repository.prs?.summary?.merge_avg + ' day(s)',
            });
          })
          setRepos(rps)

          if (start == null || end == null) {
            // first click on page or team
            const start_date = formatDate(rangeDateTime[0]);
            const end_date = formatDate(rangeDateTime[1]);

            setLoadingState(false)

            history.push(
              '/home/github?team=' +
              currentTeam +
              '&start=' +
              start_date +
              '&end=' +
              end_date
            );
          } else {
            setRangeDateTime([new Date(start), new Date(end)]);

            history.push(
              '/home/github?team=' + currentTeam +
              '&start=' + start +
              '&end=' + end
            );
          }
          setLoadingState(false)
        }
      });
    }
  }, [setRepos, currentTeam]);

  function handleChange(event, from, to) {
    setRangeDateTime([from, to]);
    params.set('start', formatDate(from));
    params.set('end', formatDate(to));
    history.push(window.location.pathname + '?' + params.toString());
  }

  const start = rangeDateTime[0];
  const end = rangeDateTime[1];

  return (
    <ModalContext.Provider value={defaultModalContext}>
      <FormModal></FormModal>
      <React.Fragment>
        {/* page title bar */}
        <Header info="Analyze the GitHub metrics overview of all your team's repositories."></Header>
        <PageSection variant={PageSectionVariants.light}>
          <Title headingLevel="h3" size={TitleSizes['2xl']}>
            GitHub metrics
            <Button
              onClick={() => navigator.clipboard.writeText(window.location.href)}
              variant="link"
              icon={<CopyIcon />}
              iconPosition="right"
            >
              Copy link
            </Button>
            <Button
              style={{ float: 'right' }}
              variant={ButtonVariant.secondary}
              onClick={modalContext.handleModalToggle}
            >
              <PlusIcon /> &nbsp; Add a repository
            </Button>
          </Title>
        </PageSection>
        {/* main content  */}
        <PageSection>
          {/* the following toolbar will contain the form (dropdowns and button) to request data to the server */}
          <Grid hasGutter>
            {/* this section will show statistics and details about GitHub metric */}
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
                          Repositories Overview
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
                      <ComposableTable repos={repos} modal={modalContext}></ComposableTable>
                    </CardBody>
                  </Card>
                </GridItem>
              )
            }
          </Grid>
          {isInvalid && !loadingState && (
            <EmptyState variant={EmptyStateVariant.xl}>
              <EmptyStateIcon icon={ExclamationCircleIcon} />
              <Title headingLevel="h1" size="lg">
                Something went wrong.
              </Title>
            </EmptyState>
          )}

        </PageSection>
      </React.Fragment>
    </ModalContext.Provider>
  );
};

export { GitHub };
