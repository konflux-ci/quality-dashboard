import React, { FC, ReactElement, useEffect, useState } from 'react';
import { ReactReduxContext, useSelector } from 'react-redux';
import {
  Title,
  PageSectionVariants,
  PageSection,
  CardHeader,
  CardTitle,
  CardBody,
  Card,
  Spinner,
  Gallery,
  Button,
  Text,
  Grid, GridItem,
  Toolbar, ToolbarItem, SearchInput
} from '@patternfly/react-core'
import { listInstalledPlugins, installPlugin } from '@app/utils/APIService';
import { Tabs, Tab } from '@patternfly/react-core';
import { useHistory } from 'react-router-dom';

function importAll(r) {
  let images = {};
   r.keys().forEach((item, index) => { images[item.replace('./', '')] = r(item); });
  return images
}

type CProps = {
  name: string,
  description: string,
  logo: any,
  category: string,
  installed: boolean,
  reason: string,
  onInstall: React.MouseEventHandler,
}

export const CardWithImageAndActions: React.FunctionComponent<CProps> = (props:CProps) => {
  const [isLoading, setIsLoading] = React.useState<boolean>(false);
  const history = useHistory();

  const onInstallClick = (e) => {
    setIsLoading(true)
    setTimeout(() => {
      props.onInstall(e)
      setIsLoading(false)
    }, 1000);
  }

  const onCardClick = ()=>{
    if(props.installed) history.push('/home/installedplugins');
  }

  return (
    <>
      <Card
        style={{minHeight: '25vh'}}
        onClick={onCardClick}
        isSelectableRaised={props.installed}
      >
        <CardHeader >
          <Grid style={{width: '100%'}}>
            <GridItem span={4}>
              <img src={props.logo.default} height={50} style={{width:'50px'}}/>
            </GridItem>
            <GridItem span={8} style={{width: '100%', textAlign: "right"}}>
              { !props.installed && props.reason == 'Available' && <Button variant="secondary" onClick={onInstallClick}> 
                  {isLoading && <span><Spinner size="md" aria-label="Installing..." />&nbsp;</span>}
                  Install 
                </Button> 
              }
              { props.reason == 'Unavailable' && <Button variant="tertiary" isDisabled ouiaId="tertiary"> Unavailable </Button> }
              { props.installed && <Button variant="tertiary" readOnly ouiaId="tertiary"> Installed </Button> }
            </GridItem>
          </Grid>
        </CardHeader>
        <CardTitle>
          <Title headingLevel='h2'>{props.name}</Title>
        </CardTitle>
        <CardBody style={{fontSize:'0.8em', maxHeight: "50%", textOverflow: "ellipsis"}}>{props.description}</CardBody>
      </Card>
    </>
  );
};

type HubProps = {
  name: string,
  id: number,
  bio?: string,
}

interface PluginResponse {
  plugin: Plugin, 
  status:{
    installed: boolean
  }
}
type Plugin = {
  name: string,
  description: string,
  logo: string,
  category: string,
  reason: string,
}

export const PHub: FC<HubProps> = ({/* destructured props */}): ReactElement => { 
  const images = importAll(require.context('../../images', false, /\.(png|jpe?g|svg)$/));
  const [activeTabKey, setActiveTabKey] = React.useState<string | number>(0);
  const [isBox, setIsBox] = React.useState<boolean>(false);
  const currentTeam = useSelector((state: any) => state.teams.Team);
  let [cards, setCards] = useState<Array<PluginResponse>>([])
  let [categories, setCategories] = useState<Array<string>>([])
  const [searchValue, setSearchValue] = React.useState('');

  const handleTabClick = (
    event: React.MouseEvent<any> | React.KeyboardEvent | MouseEvent,
    tabIndex: string | number
  ) => {
    setActiveTabKey(tabIndex);
  };

  const widths = {
    default: '100px',
    sm: '80px',
    md: '150px',
    lg: '200px',
    xl: '250px',
    '2xl': '300px'
  };

  const onChange = (value: string) => {
    setSearchValue(value);
  };

  const onFilter = (card: PluginResponse) => {
    
    if(activeTabKey != null && categories[activeTabKey] == 'all'){

      let input: RegExp;
      try {
        input = new RegExp(searchValue, 'i');
      } catch (err) {
        input = new RegExp(searchValue.replace(/[.*+?^${}()|[\]\\]/g, '\\$&'), 'i');
      }
      return card.plugin.name.search(input) >= 0;

    } else if (activeTabKey && categories[activeTabKey]){

      if (searchValue === '') {
        return true && card.plugin.category == categories[activeTabKey]
      }
  
      let input: RegExp;
      try {
        input = new RegExp(searchValue, 'i');
      } catch (err) {
        input = new RegExp(searchValue.replace(/[.*+?^${}()|[\]\\]/g, '\\$&'), 'i');
      }
      return card.plugin.name.search(input) >= 0 && categories[activeTabKey] == card.plugin.category

    }
    return true
  };

  const installTeamPlugin = (team_name: string, plugin_name:string) => {
    console.log(team_name, plugin_name)
    installPlugin(team_name, plugin_name).then(res => {
      if(res.code == 200){
        listAllPlugins();
      } else {
        throw("Error installing plugins list")
      }
    })
  };

  const listAllPlugins = () =>  {
    if(currentTeam == ''){
      console.error( "team is empty. cannot get plugins")
      return
    }
    listInstalledPlugins(currentTeam).then(res => {
      if(res.code == 200){
        setCards(res.data)
        console.log("response", res)
        let categories:string[] = ['all', ...res.data.map(item => item.plugin.category).filter((value, index, self) => self.indexOf(value) === index)]
        setCategories(categories)
      } else {
        throw("Error getting plugins list")
      }
    })
  };

  useEffect(() => {
    console.log(currentTeam)
    if(currentTeam != ''){ listAllPlugins() }
  }, [currentTeam]);

  return (
    <React.Fragment>
      <PageSection variant={PageSectionVariants.light}>
        <Title headingLevel='h1'>Plugins Hub</Title>
        <Text style={{paddingTop: "1em"}}>
          Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.
          Et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.
        </Text>
      </PageSection>
      <Grid>
          <GridItem span={3} style={{minHeight: '100%', background: 'white', padding: "1vw"}}>
            <Tabs
              style={{height: '100%'}}
              activeKey={activeTabKey}
              onSelect={handleTabClick}
              isVertical
              isBox={isBox}
              aria-label="Tabs in the vertical example"
              role="region"
            >
              { categories.map((cat, key) => (
                  <Tab style={{padding: "0.5em 1em"}} key={key} eventKey={key} title={<CardTitle style={{textTransform: 'capitalize'}}>{cat.toString()}</CardTitle>}></Tab>
                ))
              }
            </Tabs>
          </GridItem>
          <GridItem span={9} >
            <PageSection variant={PageSectionVariants.light}>
              <Toolbar id="toolbar-items-example">
                <ToolbarItem variant="search-filter"  widths={widths}>
                  <SearchInput
                    placeholder="Find by name"
                    value={searchValue}
                    onChange={(_event, value) => onChange(value)}
                    onClear={() => onChange('')}
                  />
                </ToolbarItem>
              </Toolbar>
            </PageSection>
            <Gallery hasGutter style={{margin: '1em'}}>
              {cards.filter(onFilter).map((product, key) => (
                <CardWithImageAndActions onInstall={()=>{installTeamPlugin(currentTeam, product.plugin.name)}} name={product.plugin.name} description={product.plugin.description} category={product.plugin.category} logo={images[product.plugin.logo]} reason={product.plugin.reason} installed={product.status.installed} key={key}></CardWithImageAndActions>
              ))}
            </Gallery>
        </GridItem>
      </Grid>
    </React.Fragment>
  )
};