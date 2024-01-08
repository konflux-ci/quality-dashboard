import { useSelector } from 'react-redux';

export function accessibleRouteChangeHandler() {
  return window.setTimeout(() => {
    const mainContainer = document.getElementById('primary-app-container');
    if (mainContainer) {
      mainContainer.focus();
    }
  }, 50);
}

export function teamIsNotEmpty(team: string) {
  if (team != '' && team != undefined && team != null && team != 'Select Team') {
    return true;
  }
  return false;
}

export function loadStateContext(key) {
  try {
    let serializedItem = localStorage.getItem(key);
    if (serializedItem == null) return null;
    else return JSON.parse(serializedItem);
  } catch (err) {
    return undefined;
  }
}

export function saveStateContext(key: string, item) {
  try {
    let serializedItem = JSON.stringify(item);
    localStorage.setItem(key, serializedItem);
  } catch (err) {
    console.log(err);
  }
}

export function stateContextExists(key: string) {
  try {
    if (localStorage.getItem(key) == null) return false;
    else return true;
  } catch (err) {
    return false;
  }
}

// Validates if the team pointed in the query parameter exists in the teams available
export function isValidTeam() {
  const params = new URLSearchParams(window.location.search);
  const team = params.get('team');
  const teams = useSelector((state: any) => state.teams);

  if (teams.TeamsAvailable.find((t) => t.team_name === team) || team == null) {
    return true;
  }
  return false;
}

// validateRepositoryParams validates if the 'repository' and 'organization' exists in 'repos'
export function validateRepositoryParams(repos, repository, organization) {
  if (repos.find((r) => r.organization == organization && r.repoName == repository)) {
    return true;
  }
  return false;
}

// validateParam validates if 'param' exists in 'params'
export function validateParam(params, param) {
  if (params.find((p) => p == param)) {
    return true;
  }
  return false;
}

// getLabels gets the labels, considering its name prefix
export const getLabels = (datum, prefix) => {
  if (!datum.name.startsWith(prefix)) {
    return `${datum.name} : ${datum.y}`;
  }

  return `${datum.x} \n ${datum.name} : ${datum.y}`;
};

export function sortGlobalSLI(bugs) {
  bugs.sort((a, b) => (a.global_sli > b.global_sli ? 1 : -1));
}

export const getRepoNameFormatted = (repoName) => {
  if (repoName == "infra-deployments") {
    return "RHTAP E2E Tests"
  }

  let formattedRepoName = repoName
  formattedRepoName = formattedRepoName.replaceAll("-", " ");
  formattedRepoName = formattedRepoName.toLowerCase().split(' ').map(word => word.charAt(0).toUpperCase() + word.slice(1)).join(' ');
  formattedRepoName = formattedRepoName.replace("E2e", "E2E");
  return "Individual Component - " + formattedRepoName
}