import { Component } from "react";

export function accessibleRouteChangeHandler() {
  return window.setTimeout(() => {
    const mainContainer = document.getElementById('primary-app-container');
    if (mainContainer) {
      mainContainer.focus();
    }
  }, 50);
};

export function teamIsNotEmpty(team: string) {
  if (team != "" && team != undefined && team != null && team != "Select Team") {
    return true;
  }
  return false;
};

export function loadStateContext(key) {
  try {
    let serializedItem = localStorage.getItem(key);
    if (serializedItem == null) return null;
    else return JSON.parse(serializedItem);
  } catch (err) {
    return null;
  }

}

export function saveStateContext(key: string, item) {
  try {
    let serializedItem = JSON.stringify(item);
    localStorage.setItem(key, serializedItem);
  }
  catch (err) {
    console.log(err);
  }
}

export function stateContextExists(key: string) {
  try {
    localStorage.getItem(key);
    return true;
  }
  catch(err) {
    return false;
  }
}