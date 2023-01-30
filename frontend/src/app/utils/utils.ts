import { Component } from 'react';
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
  const params = new URLSearchParams(window.location.search)
  const team = params.get("team")
  const teams = useSelector((state: any) => state.teams);

  if (teams.TeamsAvailable.find(t => t.team_name === team) || team == null) {
    return true
  }
  return false
}