import Step from '@material-ui/core/Step';
import StepLabel from '@material-ui/core/StepLabel';
import Stepper from '@material-ui/core/Stepper';
import React from 'react';
import DecryptExpand from './DecryptExpand';

function UnsealState() {
  return ['1st Decrypt key', '2nd Decrypt Key', '3rd Decrypt Key'];
}

export default function VaultDecryptRoot(props: any) {
  const steps = UnsealState();
  return (
    <div>
      <Stepper activeStep={props.sealStatus.progress} alternativeLabel>
        {steps.map((label) => (
          <Step key={label}>
            <StepLabel>{label}</StepLabel>
          </Step>
        ))}
      </Stepper>
      <DecryptExpand
        SubmitDecryptKey={props.SubmitDecryptKey}
        handleUnsealKeyInputChange={props.handleUnsealKeyInputChange}
        unsealProgress={props.unsealProgress}
      />
    </div>
  );
}
