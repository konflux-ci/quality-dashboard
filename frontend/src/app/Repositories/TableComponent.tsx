/* eslint-disable no-console */
import React, { useContext, useEffect, useState } from 'react';
import {
  TableComposable,
  Thead,
  Tr,
  Th,
  Tbody,
  Td,
  ActionsColumn,
  IAction,
  Caption,
  ThProps,
} from '@patternfly/react-table';
import {
  Alert, Pagination, PaginationVariant,
  Button, ButtonVariant,
  ExpandableSection,
  Toolbar, ToolbarItem, ToolbarContent,
  ToolbarFilter,
  ToolbarToggleGroup,
  ToolbarGroup,
  Select,
  SelectOption,
  SelectVariant,
  Title,
  EmptyStateVariant,
  EmptyState,
  EmptyStateIcon
} from '@patternfly/react-core';
import { deleteInApi, getRepositories } from '@app/utils/APIService';
import { ExternalLinkAltIcon, FilterIcon, PlusIcon, ExclamationCircleIcon } from '@patternfly/react-icons';
import _ from 'lodash';
import { useModalContext } from '@app/Repositories/CreateRepository';
import { Repository } from '@app/Repositories/';
import { ReactReduxContext, useSelector } from 'react-redux';
import { useHistory } from 'react-router-dom';
import { isValidTeam } from '@app/utils/utils';

export interface TableComponentProps {
  showCoverage?: boolean
  showDescription?: boolean
  showTableToolbar?: boolean
  enableFiltersOnTheseColumns?: Array<string>
}

const columnNames = {
  git_organization: 'GitHub Organization',
  repository_name: 'Repository',
  description: 'Description',
  coverageType: 'Coverage Type',
  coverage_percentage: 'Code Covered',
}

const rederCoverageEffects = (repo: Repository) => {
  const coveredFixed = repo.code_coverage.coverage_percentage
  if (coveredFixed >= 0 && coveredFixed <= 33.33) {
    return <Alert title={coveredFixed.toFixed(2) + "%"} variant="danger" isInline isPlain />
  } else if (coveredFixed >= 33.33 && coveredFixed <= 66.66) {
    return <Alert title={coveredFixed.toFixed(2) + "%"} variant="warning" isInline isPlain />
  }
  return <Alert title={coveredFixed.toFixed(2) + "%"} variant="success" isInline isPlain />
}

type IFilterItem = {
  State: boolean;
  Filters: Array<string | number>;
}

export const TableComponent = ({ showCoverage, showDescription, showTableToolbar, enableFiltersOnTheseColumns }: TableComponentProps) => {

  const { store } = useContext(ReactReduxContext);
  const state = store.getState();
  const dispatch = store.dispatch;

  const [perpage, onperpageset] = useState(10)
  const [repos, setRepositories] = useState<any>([])
  const [page, onPageset] = useState(1)
  const [allreps, setallreps] = useState<any>(state.repos.Allrepositories);
  const modalContext = useModalContext()

  async function deleteRepository(gitOrg: string, repoName: string) {
    const data = {
      git_organization: gitOrg,
      repository_name: repoName,
    }
    try {
      await deleteInApi(data, '/api/quality/repositories/delete')
      window.location.reload();
    } catch (error) {
      console.log(error)
    }
  }

  async function editRepository(repo: Repository) {
    try {
      modalContext.handleModalToggle(true, repo)
    } catch (error) {
      console.log(error)
    }
  }

  const defaultActions = (repo: Repository): IAction[] => [
    {
      title: 'Delete Repository',
      onClick: () => deleteRepository(repo.git_organization, repo.repository_name)
    },
    {
      title: 'Edit Repository',
      onClick: () => editRepository(repo)
    },
  ];

  function onPageselect(e, page) {
    onPageset(page)
  }

  function onperpageselect(e, Perpage) {
    onperpageset(Perpage)
  }

  // Sort helpers
  const [activeSortIndex, setActiveSortIndex] = React.useState<number | null>(null);
  const [activeSortDirection, setActiveSortDirection] = React.useState<'asc' | 'desc' | null>(null);
  const getSortableRowValues = (repo: Repository): (string | number)[] => {
    const { git_organization, repository_name, code_coverage } = repo;
    return [git_organization, repository_name, code_coverage.coverage_percentage];
  };

  let sortedRepositories = repos;

  if (activeSortIndex !== null) {
    sortedRepositories = repos.sort((a, b) => {
      const aValue = getSortableRowValues(a)[activeSortIndex];
      const bValue = getSortableRowValues(b)[activeSortIndex];
      if (typeof aValue === 'number') {
        // Numeric sort
        if (activeSortDirection === 'asc') {
          return (aValue as number) - (bValue as number);
        }
        return (bValue as number) - (aValue as number);
      } else {
        // String sort
        if (activeSortDirection === 'asc') {
          return (aValue as string).localeCompare(bValue as string);
        }
        return (bValue as string).localeCompare(aValue as string);
      }
    });
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

  // Filters helpers

  const InitFilters = (): Record<string, IFilterItem> | undefined => {
    const newObj: Record<string, IFilterItem> = {}
    if (enableFiltersOnTheseColumns != undefined) {
      for (const column of enableFiltersOnTheseColumns) {
        newObj[column] = { State: false, Filters: [] }
      }
    }
    return newObj
  };

  const Filters = InitFilters()

  const [filters, filtersDispatch] = React.useReducer(
    // filters logic
    (state, action) => {
      switch (action.type) {
        case "TOGGLE":
          state[action.payload.type].State = action.payload.isOpen
          return {
            ...state
          };

        case "SELECTION":
          if (state[action.payload.type].Filters.indexOf(action.payload.selection) === -1) {
            state[action.payload.type].Filters = [
              ...state[action.payload.type].Filters,
              action.payload.selection
            ];
          }
          return { ...state };

        case "DELETE":
          state[action.payload.type].Filters = state[action.payload.type].Filters.filter(function (value) {
            return value !== action.payload.selection;
          });
          return { ...state };

        case "DELETE_ALL":
          state[action.payload.type].Filters = [];
          return { ...state };

        case "CLEAR_ALL":
          Object.keys(state).forEach(key => {
            state[key].Filters = []
            state[key].State = false
          });
          return { ...state };

        default:
          return state;
      }
    },
    Filters
  )

  const onDelete = (category: string, chip: string) => {
    filtersDispatch({
      type: "DELETE",
      payload: {
        selection: chip,
        type: category.toLocaleLowerCase()
      }
    });
  };

  const onClearAll = () => {
    filtersDispatch({
      type: "CLEAR_ALL",
      payload: {}
    });
  };

  const onDeleteGroup = (type: string) => {
    filtersDispatch({
      type: "DELETE_ALL",
      payload: {
        type: type.toLocaleLowerCase()
      }
    });
  };

  const onToggle = (isOpen: boolean, type: string) => {
    filtersDispatch({
      type: "TOGGLE",
      payload: {
        isOpen: isOpen,
        type: type.toLocaleLowerCase()
      }
    });
  };

  const onSelect = (event, selection, type) => {
    if (!event.target.checked) {
      filtersDispatch({
        type: "DELETE",
        payload: {
          selection: selection,
          type: type.toLocaleLowerCase()
        }
      });
    } else {
      filtersDispatch({
        type: "SELECTION",
        payload: {
          selection: selection,
          type: type.toLocaleLowerCase()
        }
      });
    }
  };

  const filterData = () => {
    const filteredRepos = repos.filter(function (record) {
      let isFiltered = true;
      Object.keys(filters).forEach(category => {
        if (filters[category].Filters.length != 0) {
          const val = category.toString().split('.').reduce((o, i) => o[i], record)
          isFiltered = filters[category].Filters.includes(val.toString()) && isFiltered
        }
      });

      return isFiltered
    });

    return filteredRepos
  }

  const filteringIsActive = () => {
    let filterActive = false
    Object.keys(filters).forEach(category => { filterActive = filterActive || filters[category].Filters.length !== 0 })
    return filterActive
  }

  if (filteringIsActive()) {
    sortedRepositories = filterData()
  }

  const toggleGroupItems = (

    <ToolbarGroup variant="filter-group">
      {enableFiltersOnTheseColumns != undefined && enableFiltersOnTheseColumns.map((filter, f_idx) => {
        const f = filter.split('.').pop() || filter;
        return <ToolbarFilter
          key={"filter" + f + f_idx}
          chips={filters[filter].Filters}
          deleteChip={(category, chip) =>
            onDelete(category as string, chip as string)
          }
          deleteChipGroup={(category) => onDeleteGroup(category as string)}
          categoryName={filter}
        >
          <Select
            variant={SelectVariant.checkbox}
            aria-label="Select Input"
            onToggle={(isOpen) => onToggle(isOpen, filter)}
            onSelect={(event, selection) =>
              onSelect(event, selection, filter)
            }
            selections={filters[filter].Filters}
            isCheckboxSelectionBadgeHidden
            isOpen={filters[filter].State}
            placeholderText={"Filter by " + (columnNames[f] || filter)}
            aria-labelledby={"checkbox-select-id-" + filter}
          >
            {
              repos
                .map((value, index) => {
                  return filter.split('.').reduce((o, i) => o[i], value).toString()
                })
                .filter((x, i, a) => a.indexOf(x) == i)
                .sort((one, two) => (one < two ? -1 : 1))
                .map((v, i) => {
                  return <SelectOption key={v} value={v} />
                })
            }
          </Select>
        </ToolbarFilter>
      })}
    </ToolbarGroup>
  );

  const toolbarItems = (
    <React.Fragment>
      <ToolbarToggleGroup toggleIcon={<FilterIcon />} breakpoint="xl">
        {toggleGroupItems}
      </ToolbarToggleGroup>
    </React.Fragment>
  );
  // End of filters helpers

  const history = useHistory();
  const currentTeam = useSelector((state: any) => state.teams.Team);

  useEffect(() => {
    getRepositories(perpage, state.teams.Team).then((res) => {
      if (res.code === 200) {
        const result = res.data;
        setallreps(res.all)
        const repositories: any = result
        const repos = repositories[page - 1]
        const ress: any = []
        _.each(repos, function (ele, index, array) {
          ress.push(ele)
        })
        setRepositories(ress)
        onPageset(page)
        dispatch({ type: "SET_REPOSITORIES", data: result });
        dispatch({ type: "SET_REPOSITORIES_ALL", data: res.all });
        const params = new URLSearchParams(window.location.search)
        const team = params.get("team")
        if (team == state.teams.Team || team == null) {
          history.push('/home/repositories?team=' + currentTeam)
        }
      } else {
        dispatch({ type: "SET_ERROR", data: res });
      }
    })
  }, [page, perpage, setRepositories, dispatch, state.teams.Team, currentTeam])

  return (
    <React.Fragment>
      {isValidTeam() && <TableComposable aria-label="Actions table">
        <Caption>Repositories Summary</Caption>
        <Thead>
          <Tr>
            <Th sort={getSortParams(0)}>{columnNames.git_organization}</Th>
            <Th sort={getSortParams(1)}>{columnNames.repository_name}</Th>
            {showCoverage &&
              <Th>{columnNames.coverageType}</Th>
            }
            {showCoverage &&
              <Th sort={getSortParams(2)}>{columnNames.coverage_percentage}</Th>
            }
            {showDescription &&
              <Th>{columnNames.description}</Th>
            }
          </Tr>
        </Thead>
        <Tbody>
          {sortedRepositories.map(repo => {
            const rowActions: IAction[] | null = defaultActions(repo);
            const org_url = repo.git_url.substring(0, repo.git_url.lastIndexOf('/'));
            return (
              <Tr key={repo.repository_name}>
                <Td>
                  <a href={org_url}>
                    {repo.git_organization}
                  </a>
                  <a href={org_url}>
                    <ExternalLinkAltIcon style={{ marginLeft: "1%" }}></ExternalLinkAltIcon>
                  </a>
                </Td>
                <Td>
                  <a href={repo.git_url}>{repo.repository_name}</a>
                  <a href={repo.git_url}>
                    <ExternalLinkAltIcon style={{ marginLeft: "1%" }}></ExternalLinkAltIcon>
                  </a>

                  {!showDescription &&
                    <ExpandableSection toggleTextCollapsed='Show' toggleTextExpanded='Hide'>
                      <div>{repo.description}</div>
                    </ExpandableSection>
                  }

                </Td>
                {showCoverage &&
                  <Td><a href={`https://app.codecov.io/gh/${repo.git_organization}/${repo.repository_name}`}>CodeCov<ExternalLinkAltIcon style={{ marginLeft: "0.5%" }}></ExternalLinkAltIcon></a></Td>
                }
                {showCoverage &&
                  <Td>{rederCoverageEffects(repo)}</Td>
                }
                {showDescription &&
                  <Td>{repo.description}</Td>
                }
                <Td>
                  {rowActions ? (
                    <ActionsColumn
                      items={rowActions}
                    />
                  ) : null}
                </Td>
              </Tr>
            );
          })}
        </Tbody>
      </TableComposable>
      }
      {isValidTeam() && showTableToolbar &&
        <Toolbar id="toolbar-with-filter"
          className="pf-m-toggle-group-container"
          collapseListedFiltersBreakpoint="xl"
          clearAllFilters={onClearAll}
        >

          <ToolbarContent>
            <ToolbarItem alignment={{ default: 'alignLeft' }}>
              <Button variant={ButtonVariant.secondary} onClick={modalContext.handleModalToggle}>
                <PlusIcon /> &nbsp; Add a repository
              </Button>
            </ToolbarItem>
            <ToolbarItem>
              {toolbarItems}
            </ToolbarItem>
            <ToolbarItem alignment={{ default: 'alignRight' }}>
              <Pagination
                itemCount={allreps.length}
                perPage={perpage}
                page={page}
                variant={PaginationVariant.bottom}
                onSetPage={onPageselect}
                onPerPageSelect={onperpageselect}
              />
            </ToolbarItem>
          </ToolbarContent>
        </Toolbar>
      }
      {!isValidTeam() && <EmptyState variant={EmptyStateVariant.xl}>
        <EmptyStateIcon icon={ExclamationCircleIcon} />
        <Title headingLevel="h1" size="lg">
          Something went wrong. Please, check the URL.
        </Title>
      </EmptyState>
      }

    </React.Fragment>
  );
};