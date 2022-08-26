export function accessibleRouteChangeHandler() {
  return window.setTimeout(() => {
    const mainContainer = document.getElementById('primary-app-container');
    if (mainContainer) {
      mainContainer.focus();
    }
  }, 50);
}

export function teamIsNotEmpty(team:string){
  if(team != "" && team != undefined && team != null && team != "Select Team"){
    return true
  }
  return false
}