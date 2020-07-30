import Button from '@material-ui/core/Button';
import Dialog from '@material-ui/core/Dialog';
import DialogContent from '@material-ui/core/DialogContent';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import ListItemText from '@material-ui/core/ListItemText';
import React from 'react';

type OrgSelectProps = {
  orgs: Org[];
  submitLoginRequest: (e: React.FormEvent<Element>, idpName: string, orgID: string) => void;
};

type Org = {
  orgName: string;
  ID: string;
};

export default function OrgSelect(props: OrgSelectProps) {
  return (
    <Dialog open={props.orgs.length > 0}>
      <DialogContent>
        <List>
          {props.orgs.map((org) => (
            <Button>
              <ListItem>
                <ListItemText onClick={(e) => props.submitLoginRequest(e, 'trasa', org.ID)}>
                  {org.orgName}
                </ListItemText>
              </ListItem>
            </Button>
          ))}
        </List>
      </DialogContent>
    </Dialog>
  );
}
