import React, { useContext, useEffect, useState } from 'react';
import { CubesIcon } from '@patternfly/react-icons';
import {
  PageSection,PageSectionVariants,
  EmptyState,
  EmptyStateVariant,
  EmptyStateIcon,
  EmptyStateBody,
  DataList, DataListItem, DataListItemRow, DataListItemCells, DataListCell,
  Title, TitleSizes,
  DescriptionList,
  DescriptionListTerm,
  DescriptionListGroup,
  DescriptionListDescription,
  Alert, AlertActionLink, AlertGroup,AlertVariant, AlertActionCloseButton
} from '@patternfly/react-core';

import { Toolbar, ToolbarItem, ToolbarContent } from '@patternfly/react-core';
import { Button, SearchInput } from '@patternfly/react-core';
import { Select, SelectOption, SelectVariant } from '@patternfly/react-core';
import { getAllRepositoriesWithOrgs, getLatestProwJob } from '@app/utils/APIService';

// eslint-disable-next-line prefer-const
let Support = () => {
  const [repositories, setRepositories] = useState<{repoName: string, organization: string, isPlaceholder?: boolean}[]>([]);
  const [repoName, setRepoName] = useState("");
  const [repoOrg, setRepoOrg] = useState("");
  const [jobType, setjobType] = useState("");
  const [jobTypeToggle, setjobTypeToggle] = useState(false);
  const [repoNameToggle, setRepoNameToggle] = useState(false);
  const [prowVisible, setProwVisible] = useState(false)
  const [buttonDisabled, setbuttonDisabled] = useState(true);
  const [alerts, setAlerts] = React.useState<React.ReactNode[]>([]);
  const [prowJobs, setprowJobs] = useState([])

  const setRepoNameOnChange = (event, selection, isPlaceholder) => { 
    if (isPlaceholder){
      setRepoName("");
      setRepoOrg("");
      setRepoNameToggle(false);
    }
    else {
      setRepoName(repositories[selection].repoName); 
      setRepoOrg(repositories[selection].organization); 
      setRepoNameToggle(false) 
    }
  };

  const setjobTypeOnChange  = (event, selection, isPlaceholder) => { 
    if (isPlaceholder){
      setjobType("");
      setjobTypeToggle(false);
    }
    else{
      setjobType(selection); 
      setjobTypeToggle(false);
    }
  };

  const validateGetProwJob = () => {
    console.log(repoName, repoOrg, jobType)
    if(repoName != "" && repoOrg != "" && jobType != ""){
      setbuttonDisabled(false)
    }
    else{
      setbuttonDisabled(true)
    }
  }

  useEffect(() => {
    validateGetProwJob();
  }, [repoName, jobType]);

  const getProwJob = async () => {
    setProwVisible(true)
    try {
      let data = await getLatestProwJob(repoName, repoOrg, jobType)
      setprowJobs(data)
    }
    catch {
      setProwVisible(false);
      setAlerts(prevAlerts => {
        return [...prevAlerts,
          <Alert
          variant="danger"
          timeout={5000}
          title="Error fetching data from server"
          key={0}
          actionClose={
            <AlertActionCloseButton
              title="Error fetching data"
              variantLabel={`danger alert`}
              onClose={() => setAlerts([])}
            />
          }
        />
        ]
      });
    }
  }

  const clearProwJob = () => {
    setProwVisible(false);
    setRepoName("");
    setRepoOrg("");
    setjobType("");
  }

  const clearRepo = () => {
    setRepoName("")
    setRepoOrg("")
    setRepoNameToggle(false) 
  }

  useEffect( () => {
    getAllRepositoriesWithOrgs()
    .then((data:any) => {
      data.unshift({repoName: "Select a repository", organization: "", isPlaceholder: true})
      setRepositories(data)
    })
  }, []);

  let jobTypes = [
    <SelectOption key={0} value="periodic"/>,
    <SelectOption key={1} value="presubmit"/>,
    <SelectOption key={2} value="postsubmit"/>,
  ]

  let statusColorMap = new Map<string, string>([
    ["skipped", "lightgrey"],
    ["passed", "darkgreen"],
    ["failed", "darkred"]
  ]);

  return (
    
    <React.Fragment>
    <PageSection variant={PageSectionVariants.light}>
      <Title headingLevel="h3" size={TitleSizes['2xl']}>
        Tests Reports
      </Title>
    </PageSection>
    <PageSection>
      <AlertGroup isToast isLiveRegion> {alerts} </AlertGroup>
      <Toolbar style={{
          width:  prowVisible ? '100%' : '100%',
          margin: prowVisible ? 'auto' : '0 auto'
        }}>
          <ToolbarContent style={{textAlign: 'center'}}>
            <ToolbarItem  style={{minWidth: "20%", maxWidth: "40%"}}>
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
            <ToolbarItem style={{minWidth: "20%", maxWidth: "40%"}}>
              <Select placeholderText="Filter by status/vendor" isOpen={jobTypeToggle} onToggle={setjobTypeToggle} selections={jobType} onSelect={setjobTypeOnChange} aria-label="Select Input">
                {jobTypes}
              </Select>
            </ToolbarItem>
            <ToolbarItem >      
              <Button variant="primary" isDisabled={buttonDisabled} onClick={getProwJob}>Get Latest Test Report</Button>
            </ToolbarItem>
            <ToolbarItem >      
              <Button variant="link" onClick={clearProwJob}>Clear</Button>
            </ToolbarItem>
          </ToolbarContent>
      </Toolbar>
      
      {!prowVisible && <EmptyState variant={EmptyStateVariant.xl}>
        <EmptyStateIcon icon={CubesIcon}/>
          <Title headingLevel="h1" size="lg">
            No job selected yet.
            </Title>
          <EmptyStateBody>
            Please select a repository and an organization to see the last job's details
          </EmptyStateBody>
        </EmptyState>
      }
      <React.Fragment>
    { prowVisible && <div>
      <PageSection variant={PageSectionVariants.light} style={{marginTop: '20px'}}>
        <DescriptionList columnModifier={{default: '3Col'}}>
          <DescriptionListGroup>
            <DescriptionListTerm>Repository</DescriptionListTerm>
            <DescriptionListDescription>repo name</DescriptionListDescription>
          </DescriptionListGroup>
          <DescriptionListGroup>
            <DescriptionListTerm>Organization</DescriptionListTerm>
            <DescriptionListDescription>repo org</DescriptionListDescription>
          </DescriptionListGroup>
          <DescriptionListGroup>
            <DescriptionListTerm> % of install failures </DescriptionListTerm>
            <DescriptionListDescription> &nbsp; </DescriptionListDescription>
          </DescriptionListGroup>
          <DescriptionListGroup>
            <DescriptionListTerm> % of passed tests </DescriptionListTerm>
            <DescriptionListDescription>0</DescriptionListDescription>
          </DescriptionListGroup>
          <DescriptionListGroup>
            <DescriptionListTerm>% of failures</DescriptionListTerm>
            <DescriptionListDescription>0</DescriptionListDescription>
          </DescriptionListGroup>
          <DescriptionListGroup>
            <DescriptionListTerm>% of ci failures</DescriptionListTerm>
            <DescriptionListDescription>-</DescriptionListDescription>
          </DescriptionListGroup>
        </DescriptionList>
      </PageSection>
      <DataList aria-label="Simple data list example" style={{marginTop: '20px'}}>
        <DataListItem aria-labelledby="simple-item1">
          <DataListItemRow key="000" style={{fontWeight: 'bold', borderBottom: "1px solid lightgrey"}}>
            <DataListItemCells
              dataListCells={[
                <DataListCell width={4} key={2}>Name</DataListCell>,
                <DataListCell width={1} key={3}>Status</DataListCell>,
                <DataListCell width={1} key={4}>Time elapsed</DataListCell>,
              ]}
            />
          </DataListItemRow>
        {prowJobs.map(function(value, index){
          return <DataListItemRow key={index}>
                  <DataListItemCells
                    dataListCells={[
                      <DataListCell width={4} key={index+"-2"}>{value['name']}</DataListCell>,
                      <DataListCell width={1} key={index+"-3"} style={{fontWeight: "bold", textTransform: "uppercase", color : statusColorMap.get(value["status"])}}>
                        {value['status']}
                      </DataListCell>,
                      <DataListCell width={1} key={index+"-4"}>{value['time']}</DataListCell>,
                    ]}
                  />
                </DataListItemRow>
          })}
        </DataListItem>
      </DataList>
      </div> }
    </React.Fragment>
    </PageSection>
    </React.Fragment>
  )}

export { Support };
