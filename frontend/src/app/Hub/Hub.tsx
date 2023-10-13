import React, { FC, ReactElement, useState } from 'react';
import {
  Title,
  PageSectionVariants,
  PageSection,
  CardHeader,
  CardTitle,
  CardBody,
  Card,
  CardFooter,
  Gallery,
  Button,
  Text,
  Grid, GridItem,
  Toolbar, ToolbarItem, SearchInput
} from '@patternfly/react-core'

import { Tabs, Tab, TabTitleText, Checkbox, Tooltip } from '@patternfly/react-core';

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
  status: string,
}

export const CardWithImageAndActions: React.FunctionComponent<CProps> = (props:CProps) => {
  const [isOpen, setIsOpen] = React.useState<boolean>(false);

  return (
    <>
      <Card>
        <CardHeader>
          <img src={props.logo.default} height={50} style={{width:'50px'}}/>
        </CardHeader>
        <CardTitle>
          <Title headingLevel='h2'>{props.name}</Title>
        </CardTitle>
        <CardBody>{props.description}</CardBody>
        <CardFooter>
          { props.status == 'available' && <Button variant="primary" ouiaId="Primary"> Install </Button> }
          { props.status == 'unavailable' && <Button variant="tertiary" isDisabled ouiaId="Primary"> Unavailable </Button> }
          { props.status == 'installed' && <Button variant="secondary" readOnly ouiaId="Primary"> Installed </Button> }
        </CardFooter>
      </Card>
    </>
  );
};

type HubProps = {
  name: string,
  id: number,
  bio?: string,
}

type Plugin = {
  name: string,
  description: string,
  logo: string,
  category: string,
  status: string,
}

export const PHub: FC<HubProps> = ({/* destructured props */}): ReactElement => { 
  const images = importAll(require.context('../../images', false, /\.(png|jpe?g|svg)$/));
  const [activeTabKey, setActiveTabKey] = React.useState<string | number>(0);
  const [isBox, setIsBox] = React.useState<boolean>(false);
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
  let c:Array<Plugin> = [
    {
      name: "Github",
      description: "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. ",
      logo: "github-mark.png",
      category: "github",
      status: "available",
    },
    {
      name: "Openshift CI",
      description: "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. ",
      logo: "github-mark.png",
      category: "openshift",
      status: "unavailable",
    },
    {
      name: "Codecov",
      description: "some description here",
      logo: "codecov-square.png",
      category: "codecov",
      status: "available",
    },
    {
      name: "Some Plugin",
      description: "some description here",
      logo: "github-mark.png",
      category: "github",
      status: "installed",
    },
    {
      name: "Jira Issues",
      description: "some description here for jira yeeeeeeeeeeeeee",
      logo: "jira.png",
      category: "jira",
      status: "available",
    }
  ]
  const categories:string[] = ['all', ...c.map(item => item.category).filter((value, index, self) => self.indexOf(value) === index)]

  let [cards, setCards] = useState<Array<Plugin>>(c)
  const [searchValue, setSearchValue] = React.useState('');
  const onChange = (value: string) => {
    setSearchValue(value);
  };

  const onFilter = (card: Plugin) => {
    
    if(activeTabKey != null && categories[activeTabKey] == 'all'){

      let input: RegExp;
      try {
        input = new RegExp(searchValue, 'i');
      } catch (err) {
        input = new RegExp(searchValue.replace(/[.*+?^${}()|[\]\\]/g, '\\$&'), 'i');
      }
      return card.name.search(input) >= 0;

    } else if (activeTabKey && categories[activeTabKey]){

      if (searchValue === '') {
        return true && card.category == categories[activeTabKey]
      }
  
      let input: RegExp;
      try {
        input = new RegExp(searchValue, 'i');
      } catch (err) {
        input = new RegExp(searchValue.replace(/[.*+?^${}()|[\]\\]/g, '\\$&'), 'i');
      }
      return card.name.search(input) >= 0 && categories[activeTabKey] == card.category

    }
    return true
  };

  const filteredCards = cards.filter(onFilter);

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
              {filteredCards.map((product, key) => (
                <CardWithImageAndActions name={product.name} description={product.description} category={product.category} logo={images[product.logo]} status={product.status} key={key}></CardWithImageAndActions>
              ))}
            </Gallery>
        </GridItem>
      </Grid>
    </React.Fragment>
  )
};