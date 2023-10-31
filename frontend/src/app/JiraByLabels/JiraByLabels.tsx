import React, { useEffect, useState } from 'react';
import { CopyIcon, HelpIcon } from '@patternfly/react-icons';
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
  FormGroup,
  TextInput,
  Popover,
  Form,
  Modal,
} from '@patternfly/react-core';
import { Button } from '@patternfly/react-core';
import {
  getJiraIssuesByJQLQuery,
  getTeams,
} from '@app/utils/APIService';
import { Grid, GridItem } from '@patternfly/react-core';
import { ReactReduxContext, useSelector } from 'react-redux';
import { useHistory } from 'react-router-dom';
import { validateParam } from '@app/utils/utils';
import { Header } from '@app/utils/Header';
import { CustomStackChart } from './CustomStackChart';
import { Issue, ListIssues } from './ListIssuesTable';

export const defaultFilter = `project in (DEVHAS, SRVKP, GITOPSRVCE, HACBS, RHTAP, RHTAPBUGS) AND labels = ci-fail`
export const defaultLabels = ["ci-fail", "test_bug", "product_bug", "untriaged", "infra_bug", "to_investigate"]

// eslint-disable-next-line prefer-const
export let JiraByLabels = () => {
  const { store } = React.useContext(ReactReduxContext);
  const state = store.getState();
  const currentTeam = useSelector((state: any) => state.teams.Team);
  const history = useHistory();
  const params = new URLSearchParams(window.location.search);
  const [loadingState, setLoadingState] = useState(false);
  const [labels, setLabels] = useState<string[]>(defaultLabels);
  const [jiraFilter, setJiraFilter] = useState<string>(defaultFilter);
  const [issues, setIssues] = useState<any>({});
  const [isSelected, setIsSelected] = React.useState("");
  const [openIssuesTable, setOpenIssuesTable] = useState<Array<Issue>>([]);
  const [closedIssuesTable, setClosedIssuesTable] = useState<Array<Issue>>([]);

  useEffect(() => {
    getJiraIssuesByJQLQuery(jiraFilter).then((res: any) => {
      if (res.data.length > 0) {
        setIssues(res.data)
      } else {
        setIssues([])
      }
    })

  }, [jiraFilter]);

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
      setIssues([])

      const team = params.get('team');
      const jiraFilterParam = params.get('jira_filter');
      const labelsParam = params.get('labels');

      getJiraIssuesByJQLQuery(jiraFilter).then((res: any) => {
        if (res.data.length < 1 && (team == state.teams.Team || team == null)) {
          setLoadingState(false)
          history.push('/home/jira-by-labels?team=' + currentTeam);
        }

        if (res.data.length > 0 && (team == state.teams.Team || team == null)) {
          setIssues(res.data)

          if (jiraFilterParam == null || labelsParam == null) {
            // first click on page or team
            setIsSelected(labels[0])
            setLoadingState(false)

            history.push(
              '/home/jira-by-labels?team=' +
              currentTeam +
              '&jira_filter=' +
              jiraFilter +
              '&labels=' +
              labels
            );
          } else {
            setJiraFilter(jiraFilterParam);
            const lbls = labelsParam.split(",")
            setLabels(lbls);
            setIsSelected(lbls[0])

            history.push(
              '/home/jira-by-labels?team=' + currentTeam +
              '&jira_filter=' +
              jiraFilterParam +
              '&labels=' +
              labelsParam
            );

            setLoadingState(false)
          }
        }
      });
    }
  }, [setIssues, currentTeam]);

  useEffect(() => {
    if (issues.length > 0) {
      const issuesSelected = issues?.filter((x) => {
        if (x.labels.includes(isSelected)) {
          return x
        }
      })

      setOpenIssuesTable(issuesSelected.filter((x) => {
        if (x.status != "Closed") {
          return x
        }
      }))

      setClosedIssuesTable(issuesSelected.filter((x) => {
        if (x.status == "Closed") {
          return x
        }
      }))
    }
  }, [isSelected, issues]);


  const handleItemClick = (isSelected: boolean, event: React.MouseEvent<any> | React.KeyboardEvent | MouseEvent) => {
    const id = event.currentTarget.id;
    setOpenIssuesTable([])
    setClosedIssuesTable([])
    setIsSelected(id)
  };


  const [isModalOpen, setIsModalOpen] = React.useState(false);
  const [jiraFilterValue, setJiraFilterValue] = React.useState(defaultFilter);
  const [labelsValue, setLabelsValue] = React.useState("");
  type validate = 'success' | 'warning' | 'error' | 'default';
  const [labelsValidated, setLabelsValidated] = React.useState<validate>('error');
  const regexp = new RegExp('^[a-zA-Z_-]+(,[0-9a-zA-Z_-]+)*$')


  const handleModalToggle = () => {
    setIsModalOpen(!isModalOpen);
  };

  const submit = () => {
    const ls = labelsValue.split(",")
    setIsModalOpen(!isModalOpen);
    setJiraFilter(jiraFilterValue);
    setLabels(ls)
    setIsSelected(ls[0])

    history.push(
      '/home/jira-by-labels?team=' +
      params.get('team') +
      '&jira_filter=' +
      jiraFilterValue +
      '&labels=' +
      labelsValue
    );
  };

  const handleJiraKeyInput = async (value) => {
    setJiraFilterValue(value);
  };

  const handleLabelsInput = async (value) => {
    setLabelsValidated('error');
    setLabelsValue(value);
    if (regexp.test(value)) {
      setLabelsValidated('success');
    } else {
      setLabelsValidated('error');
    }
  };

  return (
    <React.Fragment>
      {/* page title bar */}
      <Header info="Create dynamically your own tables."></Header>
      <PageSection variant={PageSectionVariants.light}>
        <Title headingLevel="h3" size={TitleSizes['2xl']}>
          Jira By Labels
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


          {/* this section will show statistics and details about GitHub metric */}
          {loadingState && <div style={{ width: '100%', textAlign: "center" }}>
            <Spinner isSVG diameter="80px" aria-label="Contents of the custom size example" style={{ margin: "100px auto" }} />
          </div>
          }
          {!loadingState && (issues.length != undefined) &&
            (
              <Grid hasGutter>
                <Button variant="primary" onClick={handleModalToggle}>
                  Set configuration
                </Button>
                <Modal
                  width={800}
                  title="Configure your own tables"
                  isOpen={isModalOpen}
                  onClose={handleModalToggle}
                  actions={[
                    <Button key="confirm" variant="primary" onClick={submit} isDisabled={labelsValidated == "error"}>
                      Confirm
                    </Button>,
                    <Button key="cancel" variant="link" onClick={handleModalToggle}>
                      Cancel
                    </Button>
                  ]}
                >
                  <Form isHorizontal id="modal-with-form-form">
                    <FormGroup
                      label="Jira Filter"
                      labelIcon={
                        <Popover headerContent={<div></div>} bodyContent={<div>Add a valid Jira Key that refers to a ci-fail issue</div>}>
                          <button
                            type="button"
                            aria-label="More info for name field"
                            onClick={(e) => e.preventDefault()}
                            aria-describedby="modal-with-form-form-name"
                            className="pf-c-form__group-label-help"
                          >
                            <HelpIcon noVerticalAlign />
                          </button>
                        </Popover>
                      }
                      isRequired
                      fieldId="modal-with-form-form-name"
                      helperTextInvalid="Must be a valid JIRA key"
                    >
                      <TextInput
                        isRequired
                        type="email"
                        id="modal-with-form-form-name"
                        name="modal-with-form-form-name"
                        value={jiraFilterValue}
                        onChange={handleJiraKeyInput}
                      />
                    </FormGroup>
                    <FormGroup
                      label="Labels"
                      labelIcon={
                        <Popover headerContent={<div></div>} bodyContent={<div>Add a list of labels separated by comma. Example: test_bug,product_bug,to_investigate</div>}>
                          <button
                            type="button"
                            aria-label="More info for name field"
                            onClick={(e) => e.preventDefault()}
                            aria-describedby="modal-with-form-form-name"
                            className="pf-c-form__group-label-help"
                          >
                            <HelpIcon noVerticalAlign />
                          </button>
                        </Popover>
                      }
                      isRequired
                      fieldId="modal-with-form-form-name"
                      helperTextInvalid="Must be a valid JIRA key"
                    >
                      <TextInput
                        validated={labelsValidated}
                        isRequired
                        type="email"
                        id="modal-with-form-form-name"
                        name="modal-with-form-form-name"
                        value={labelsValue}
                        onChange={handleLabelsInput}
                      />
                    </FormGroup>
                  </Form>
                </Modal>

                {(openIssuesTable.length > 0) &&
                  <GridItem span={4} rows={12}>
                    <Card style={{ textAlign: 'center', alignContent: 'center' }}>
                      <CardTitle style={{ backgroundColor: "grey", color: 'white' }}>Open Issues by Labels</CardTitle>
                    </Card>
                    <CardBody style={{ backgroundColor: 'white' }}>
                      <CustomStackChart data={[getLabels(openIssuesTable, labels)]} legend={[{ name: 'Issues' }]} />
                    </CardBody>
                  </GridItem>
                }
                {(openIssuesTable.length > 0) &&
                  <GridItem span={4} rows={12}>
                    <Card style={{ textAlign: 'center', alignContent: 'center' }}>
                      <CardTitle style={{ backgroundColor: "grey", color: 'white' }}>Open Issues by Labels and Component</CardTitle>
                    </Card>
                    <CardBody style={{ backgroundColor: 'white' }}>
                      <CustomStackChart data={getIssuesByFields(openIssuesTable, labels, "component")} legend={[{ name: 'ci-fail' }, { name: 'test_bug' }, { name: 'product_bug' }, { name: 'untriaged' }, { name: 'infra_bug' }, { name: 'to_investigate' }]} />
                    </CardBody>
                  </GridItem>
                }
                {(openIssuesTable.length > 0) &&
                  <GridItem span={4} rows={12}>
                    <Card style={{ textAlign: 'center', alignContent: 'center' }}>
                      <CardTitle style={{ backgroundColor: "grey", color: 'white' }}>Open Issues by Labels and Status</CardTitle>
                    </Card>
                    <CardBody style={{ backgroundColor: 'white' }} >
                      <CustomStackChart data={getIssuesByFields(openIssuesTable, labels, "status")} legend={[{ name: 'ci-fail' }, { name: 'test_bug' }, { name: 'product_bug' }, { name: 'untriaged' }, { name: 'infra_bug' }, { name: 'to_investigate' }]} />
                    </CardBody>
                  </GridItem>
                }
                {(closedIssuesTable.length > 0) &&
                  <GridItem span={4} rows={12}>
                    <Card style={{ textAlign: 'center', alignContent: 'center' }}>
                      <CardTitle style={{ backgroundColor: "grey", color: 'white' }}>Closed Issues by Labels</CardTitle>
                    </Card>
                    <CardBody style={{ backgroundColor: 'white' }}>
                      <CustomStackChart data={[getLabels(closedIssuesTable, labels)]} legend={[{ name: 'Issues' }]} />
                    </CardBody>
                  </GridItem>
                }
                {(closedIssuesTable.length > 0) &&
                  <GridItem span={4} rows={12}>
                    <Card style={{ textAlign: 'center', alignContent: 'center' }}>
                      <CardTitle style={{ backgroundColor: "grey", color: 'white' }}>Closed Issues by Labels and Component</CardTitle>
                    </Card>
                    <CardBody style={{ backgroundColor: 'white' }}>
                      <CustomStackChart data={getIssuesByFields(closedIssuesTable, labels, "component")} legend={[{ name: 'ci-fail' }, { name: 'test_bug' }, { name: 'product_bug' }, { name: 'untriaged' }, { name: 'infra_bug' }, { name: 'to_investigate' }]} />
                    </CardBody>
                  </GridItem>
                }
                {(closedIssuesTable.length > 0) &&
                  <GridItem span={4} rows={12}>
                    <Card style={{ textAlign: 'center', alignContent: 'center' }}>
                      <CardTitle style={{ backgroundColor: "grey", color: 'white' }}>Closed Issues by Labels and Status</CardTitle>
                    </Card>
                    <CardBody style={{ backgroundColor: 'white' }} >
                      <CustomStackChart data={getIssuesByFields(closedIssuesTable, labels, "status")} legend={[{ name: 'ci-fail' }, { name: 'test_bug' }, { name: 'product_bug' }, { name: 'untriaged' }, { name: 'infra_bug' }, { name: 'to_investigate' }]} />
                    </CardBody>
                  </GridItem>
                }
                {(openIssuesTable.length > 0) &&
                  <Card>
                    <CardTitle style={{ backgroundColor: "grey", color: 'white' }}>
                      List of open issues
                    </CardTitle>
                    <CardBody style={{ backgroundColor: 'white' }}>
                      <div style={{ marginTop: 10 }}>
                        <ToggleGroup aria-label="Default with single selectable">
                          {labels?.map((label, idx) => {
                            return (
                              <ToggleGroupItem
                                key={idx}
                                text={label}
                                buttonId={label}
                                isSelected={isSelected === label}
                                onChange={handleItemClick}
                              />
                            )
                          })}
                        </ToggleGroup>
                      </div>
                      <ListIssues issues={openIssuesTable}></ListIssues>
                    </CardBody>
                  </Card>
                }
                {(closedIssuesTable.length > 0) &&
                  <Card>
                    <CardTitle style={{ backgroundColor: "grey", color: 'white' }}>
                      List of closed issues
                    </CardTitle>
                    <CardBody style={{ backgroundColor: 'white' }}>
                      <div style={{ marginTop: 10 }}>
                        <ToggleGroup aria-label="Default with single selectable">
                          {labels?.map((label, idx) => {
                            return (
                              <ToggleGroupItem
                                key={idx}
                                text={label}
                                buttonId={label}
                                isSelected={isSelected === label}
                                onChange={handleItemClick}
                              />
                            )
                          })}
                        </ToggleGroup>
                      </div>
                      <ListIssues issues={closedIssuesTable}></ListIssues>
                    </CardBody>
                  </Card>
                }
              </Grid>
            )
          }
        </Grid>
      </PageSection>
    </React.Fragment>
  );
};

const getLabels = (issues: Issue[], labels: string[]) => {
  const issuesByLabels = labels?.map((x) => {
    const count = issues?.filter(y => y["labels"].includes(x) == true).length
    return {
      name: "Issues",
      x: x,
      y: count,
      label: "Issues: " + count,
    };
  }).filter(
    (elem, index, arr) => index === arr.findIndex((t) => t.x === elem.x)
  );

  issuesByLabels.sort((a, b) => (a.y < b.y ? -1 : 1));

  return issuesByLabels
}

const getIssuesByField = (issues: Issue[], label: string, field: string) => {
  const issuesByField = issues?.map((x) => {
    const count = issues?.filter(y => y[field] == x[field] && y["labels"].includes(label) == true).length

    return {
      name: label,
      x: x[field],
      y: count,
      label: label + ': ' + count
    };
  }).filter(
    (elem, index, arr) => index === arr.findIndex((t) => t.x === elem.x)
  );

  return issuesByField
}

interface Data {
  name: string,
  x: string,
  y: number,
  label: string,
}

const getIssuesByFields = (issues: Issue[], labels: string[], field: string) => {
  const issuesByFields: Data[][] = []

  labels.forEach((label) => {
    const issuesByField = getIssuesByField(issues, label, field)

    issuesByFields.push(issuesByField)

  });

  return issuesByFields
}