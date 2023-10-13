import React, { FC, ReactElement, useState } from 'react';
import {
  Title,
  PageSectionVariants,
  PageSection,
  CardHeader,
  CardTitle,
  CardBody,
  Card,
  Brand,
  CardFooter,
  Checkbox,
  Dropdown,
  DropdownItem,
  MenuToggle,
  MenuToggleElement,
  Divider,
  Gallery,
  Button
} from '@patternfly/react-core'

type CProps = {
  name: string,
  description: string,
  logo: string,
  status: 'installed' | 'available' | 'unavailable',
}

export const CardWithImageAndActions: React.FunctionComponent<CProps> = (props:CProps) => {
  const [isOpen, setIsOpen] = React.useState<boolean>(false);
  const [isChecked, setIsChecked] = React.useState<boolean>(false);
  const [hasNoOffset, setHasNoOffset] = React.useState<boolean>(false);

  const onSelect = () => {
    setIsOpen(!isOpen);
  };
  const onClick = (checked: boolean) => {
    setIsChecked(checked);
  };
  const toggleOffset = (checked: boolean) => {
    setHasNoOffset(checked);
  };

  const dropdownItems = (
    <>
      <DropdownItem key="action">Action</DropdownItem>
      {/* Prevent default onClick functionality for example purposes */}
      <DropdownItem key="link" onClick={(event: any) => event.preventDefault()}>
        Link
      </DropdownItem>
      <DropdownItem key="disabled action" isDisabled>
        Disabled Action
      </DropdownItem>
      <DropdownItem key="disabled link" isDisabled onClick={(event: any) => event.preventDefault()}>
        Disabled Link
      </DropdownItem>
      <Divider component="li" key="separator" />
      <DropdownItem key="separated action">Separated Action</DropdownItem>
      <DropdownItem key="separated link" onClick={(event: any) => event.preventDefault()}>
        Separated Link
      </DropdownItem>
    </>
  );

  return (
    <>
      <Card>
        <CardHeader>
          <img src={"/images/"+props.logo} alt="PatternFly logo" style={{ width: '50px' }} />
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

export const PHub: FC<HubProps> = ({/* destructured props */}): ReactElement => { 
  let c = [
    {
      name: "Github",
      description: "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. ",
      logo: "github-mark.png",
      status: "available",
    },
    {
      name: "Openshift CI",
      description: "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. ",
      logo: "github-mark.png",
      status: "unavailable",
    },
    {
      name: "Codecov",
      description: "some description here",
      logo: "codecov-square.png",
      status: "installed",
    }
  ]
  let [cards, setCards] = useState<Array<any>>(c)

  return (
    <React.Fragment>
      <PageSection variant={PageSectionVariants.light}>
        <Title headingLevel='h1'>Plugins Hub</Title>
        <Title headingLevel='h4'>Plugins Hub some subtitle goes here :troll:</Title>
      </PageSection>
      <PageSection variant={PageSectionVariants.light}>
        <Title headingLevel='h4'>Some searchbar here</Title>
      </PageSection>
      <PageSection variant={PageSectionVariants.default}>
      <Gallery hasGutter>
        {cards.map((product, key) => (
          <CardWithImageAndActions name={product.name} description={product.description} logo={product.logo} status={product.status} key={key}></CardWithImageAndActions>
        ))}
      </Gallery>
      </PageSection>
    </React.Fragment>
  )
};