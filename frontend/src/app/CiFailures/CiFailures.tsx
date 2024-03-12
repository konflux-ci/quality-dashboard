import React, { useContext, useEffect, useState } from 'react';
import { CopyIcon, OpenDrawerRightIcon, PlusIcon } from '@patternfly/react-icons';
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
  Flex,
  FlexItem,
  Drawer,
  DrawerContent,
  DrawerPanelContent,
  DrawerHead,
  DrawerActions,
  DrawerCloseButton,
  TextContent,
} from '@patternfly/react-core';
import { Button } from '@patternfly/react-core';
import { Grid, GridItem } from '@patternfly/react-core';
import { ComposableTable } from './Table';
import { FormModal, ModalContext, useDefaultModalContextState, useModalContext } from './CreateFailure';
import { ReactReduxContext, useSelector } from 'react-redux';
import { getFailures, getTeams, listUsers } from '@app/utils/APIService';
import { validateParam } from '@app/utils/utils';
import { formatDateTime, getRangeDates } from '@app/Reports/utils';
import { useHistory } from 'react-router-dom';
import { DateTimeRangePicker } from '@app/utils/DateTimeRangePicker';
import { Header } from '@app/utils/Header';

export interface FailureInfo {
  jira_key: string;
  jira_status: string;
  error_message: string;
  frequency: string;
  title_from_jira: string;
  created_date: string;
  closed_date: string;
  labels: string;
}

// eslint-disable-next-line prefer-const
let CiFailures = () => {
  const [loadingState, setLoadingState] = useState(false);
  const defaultModalContext = useDefaultModalContextState();
  const modalContext = useModalContext();
  const { store } = useContext(ReactReduxContext);
  const state = store.getState();
  const currentTeam = useSelector((state: any) => state.teams.Team);
  const params = new URLSearchParams(window.location.search);
  const [rangeDateTime, setRangeDateTime] = useState(getRangeDates(15));
  const [failures, setFailures] = useState<any>({});
  const history = useHistory();
  const [isExpanded, setIsExpanded] = React.useState(false);
  const drawerRef = React.useRef<HTMLDivElement>();

  function handleChange(event, from, to) {
    setRangeDateTime([from, to]);
    params.set('start', formatDateTime(from));
    params.set('end', formatDateTime(to));
    history.push(window.location.pathname + '?' + params.toString());
  }

  useEffect(() => {
    if (state.teams.Team != '') {
      getFailures(state.teams.Team, rangeDateTime).then((res: any) => {
        setFailures(res.data)
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
      setFailures([])

      const team = params.get('team');
      const start = params.get('start');
      const end = params.get('end');

      getFailures(state.teams.Team, rangeDateTime).then((res: any) => {
        if (res.data.length < 1 && (team == state.teams.Team || team == null)) {
          setLoadingState(false)
          history.push('/home/rhtapbugs-impact?team=' + currentTeam);
        }

        if (res.data.length > 0 && (team == state.teams.Team || team == null)) {
          if (start == null || end == null) {
            // first click on page or team
            const start_date = formatDateTime(rangeDateTime[0]);
            const end_date = formatDateTime(rangeDateTime[1]);

            setFailures(res.data)
            setLoadingState(false)

            history.push(
              '/home/rhtapbugs-impact?team=' +
              currentTeam +
              '&start=' +
              start_date +
              '&end=' +
              end_date
            );
          } else {
            setRangeDateTime([new Date(start), new Date(end)]);

            // getFailures(state.teams.Team, rangeDateTime).then((data: any) => {
            //   setFailures(data)
            // });

            history.push(
              '/home/rhtapbugs-impact?team=' + currentTeam +
              '&start=' + start +
              '&end=' + end
            );
          }
          setLoadingState(false)
        }
      });
    }
  }, [setFailures, currentTeam]);


  const start = rangeDateTime[0];
  const end = rangeDateTime[1];

  const onClick = () => {
    modalContext.handleModalToggle()
  }

  const onCloseClick = () => {
    setIsExpanded(false);
  };

  const onLearnMoreClick = () => {
    setIsExpanded(!isExpanded);
  }

  const onExpand = () => {
    drawerRef.current && drawerRef.current.focus();
  };

  const panelContent = (
    <DrawerPanelContent isResizable defaultSize={'500px'} minSize={'150px'}>
      <DrawerHead>
        <TextContent>
          <Title headingLevel="h1">Bug CI Impact</Title>
          <span>
            This page aims to help you observe the impact of the bugs that are affecting your team's Openshift CI prow jobs. <br />You can add, update, or delete them. To add a new entry, you need to point out the Jira Key of the bug and the associated error message.
          </span>
          <Title headingLevel="h1">How the frequency/impact is measured?</Title>
          <span>
            To calculate the impact of each bug, in the date time range selected, Quality Dashboard will search for all the team's OpenShift CI prow jobs and verify in how many of them the bug's error message is present.
            <br /> Note that the only the OpenShift CI prow jobs reports that matches regex &apos;(j?unit|e2e)-?[0-9a-z]+\.xml&apos; are saved on the DB.
          </span>
        </TextContent>
        <DrawerActions>
          <DrawerCloseButton onClick={onCloseClick} />
        </DrawerActions>
      </DrawerHead>
    </DrawerPanelContent>
  );

  return (
    <ModalContext.Provider value={defaultModalContext}>
      <React.Fragment>
        {/* page title bar */}
        <Header info="Observe the impact of the bugs that are affecting your team's Openshift CI prow jobs."></Header>
        <PageSection variant={PageSectionVariants.light}>
          <Title headingLevel="h3" size={TitleSizes['2xl']}>
            Bug CI Impact
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
              onClick={onClick}
            >
              <PlusIcon /> &nbsp; Add a bug
            </Button>
          </Title>
        </PageSection>
        {/* main content  */}
        <Drawer isExpanded={isExpanded} onExpand={onExpand}>
          <DrawerContent panelContent={panelContent}>
            <PageSection>
              {/* the following toolbar will contain the form (dropdowns and button) to request data to the server */}
              <Grid hasGutter>
                <FormModal></FormModal>


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
                              Bug CI Impact Overview
                              <Button onClick={onLearnMoreClick} variant="link" icon={<OpenDrawerRightIcon />} iconPosition="right">
                                Learn more
                              </Button>
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
                          <ComposableTable failures={failures} modal={modalContext}></ComposableTable>
                        </CardBody>
                      </Card>
                    </GridItem>
                  )
                }
              </Grid>
              {/* {isInvalid && !loadingState && (
            <EmptyState variant={EmptyStateVariant.xl}>
              <EmptyStateIcon icon={ExclamationCircleIcon} />
              <Title headingLevel="h1" size="lg">
                Something went wrong.
              </Title>
            </EmptyState>
          )} */}

            </PageSection>
          </DrawerContent>
        </Drawer>
      </React.Fragment>
    </ModalContext.Provider>
  );
};

export { CiFailures };
