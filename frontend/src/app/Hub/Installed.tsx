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
  Icon,
  Button,
  Text,
  Grid, GridItem,
  Toolbar, ToolbarItem, SearchInput
} from '@patternfly/react-core'
import { TableComposable, Thead, Tr, Th, Tbody, Td } from '@patternfly/react-table';
import CheckCircleIcon from '@patternfly/react-icons/dist/esm/icons/check-circle-icon';

function importAll(r) {
  let images = {};
   r.keys().forEach((item, index) => { images[item.replace('./', '')] = r(item); });
  return images
}

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

const columnNames = {
  name: "Name",
  description: "Description",
  logo: "",
  category: "Category",
  status: "Status",
};

export const PInstalled: FC<HubProps> = ({/* destructured props */}): ReactElement => { 
  const images = importAll(require.context('../../images', false, /\.(png|jpe?g|svg)$/));

  const widths = {
    default: '100px',
    sm: '80px',
    md: '150px',
    lg: '300px',
    xl: '350px',
    '2xl': '500px'
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
      status: "installed",
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

  let [cards, setCards] = useState<Array<Plugin>>(c)
  const [searchValue, setSearchValue] = React.useState('');
  const onChange = (value: string) => {
    setSearchValue(value);
  };

  const onFilter = (card: Plugin) => {
    
    if (searchValue === '') {
      return true
    }

    let input: RegExp;
    try {
      input = new RegExp(searchValue, 'i');
    } catch (err) {
      input = new RegExp(searchValue.replace(/[.*+?^${}()|[\]\\]/g, '\\$&'), 'i');
    }
    return card.name.search(input) >= 0 

  };

  const filteredCards = cards.filter(onFilter);

  return (
    <React.Fragment>
      <PageSection variant={PageSectionVariants.light}>
        <Title headingLevel='h1'>Installed Plugins</Title>
        <Text style={{paddingTop: "1em"}}>
          Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.
          Et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.
        </Text>
      </PageSection>
      <Grid>
          <GridItem span={12} >
            <PageSection variant={PageSectionVariants.light}>
              <Toolbar id="toolbar-items-example">
                <ToolbarItem variant="search-filter" widths={widths}>
                  <SearchInput
                    placeholder="Find by name"
                    value={searchValue}
                    onChange={(_event, value) => onChange(value)}
                    onClear={() => onChange('')}
                  />
                </ToolbarItem>
              </Toolbar>
            </PageSection>
            <PageSection variant={PageSectionVariants.light}>
              <TableComposable aria-label="Sortable table custom toolbar">
                <Thead>
                  <Tr>
                    <Th></Th>
                    <Th>{columnNames.name}</Th>
                    <Th width={70}>{columnNames.description}</Th>
                    <Th>{columnNames.category}</Th>
                    <Th>{columnNames.status}</Th>
                  </Tr>
                </Thead>
                <Tbody>
                  {filteredCards.filter((p, i)=>{return p.status=='installed'}).map((product, rowIndex) => (
                    <Tr key={rowIndex}>
                      <Td style={{width:'50px'}}><img src={images[product.logo].default} width={50} style={{width:'50px'}}/></Td>
                      <Td style={{verticalAlign:'top'}} dataLabel={columnNames.name}>{product.name}</Td>
                      <Td style={{verticalAlign:'top'}} dataLabel={columnNames.description}>{product.description}</Td>
                      <Td style={{verticalAlign:'top', textTransform: 'capitalize'}} dataLabel={columnNames.category}>{product.category}</Td>
                      <Td style={{verticalAlign:'top', textTransform: 'capitalize'}} dataLabel={columnNames.status}><Icon status="success"><CheckCircleIcon /></Icon> &nbsp; {product.status}</Td>
                    </Tr>
                  ))}
                </Tbody>
              </TableComposable>
            </PageSection>
        </GridItem>
      </Grid>
    </React.Fragment>
  )
};