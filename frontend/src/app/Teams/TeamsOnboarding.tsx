import React, {useState} from 'react';
import { Wizard, PageSection, PageSectionVariants } from '@patternfly/react-core';

export const TeamsWizard = () => {
  
  const [ stepIdReached, setState ] = useState(1);

  const onNext = (id) => {
    setState(stepIdReached < id ? id : stepIdReached);
  };

  const closeWizard = () => {
    console.log('close wizard');
  };

  const steps = [
    { id: 'incrementally-enabled-1', name: 'First step', component: <p>Step 1 content</p> },
    { id: 'incrementally-enabled-2', name: 'Second step', component: <p>Step 2 content</p>},
    { id: 'incrementally-enabled-3', name: 'Third step', component: <p>Step 3 content</p>},
    { id: 'incrementally-enabled-4', name: 'Fourth step', component: <p>Step 4 content</p>},
    {
      id: 'incrementally-enabled-5',
      name: 'Review',
      component: <p>Review step content</p>,
      nextButtonText: 'Finish'
    }
  ];
  const title = 'Incrementally enabled wizard';

  return (
    <React.Fragment>
      <PageSection style={{backgroundColor: 'white'}} variant={PageSectionVariants.light}>
        <Wizard
          navAriaLabel={`${title} steps`}
          mainAriaLabel={`${title} content`}
          onClose={closeWizard}
          steps={steps}
          onNext={onNext}
          height={600}
        />
      </PageSection>
    </React.Fragment>
  );
}