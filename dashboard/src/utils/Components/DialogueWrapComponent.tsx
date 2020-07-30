import Button from '@material-ui/core/Button';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogTitle from '@material-ui/core/DialogTitle';
import Typography from '@material-ui/core/Typography';
import React from 'react';

type DialogueProps = {
  open: boolean;
  fullScreen: boolean;
  handleClose: () => void;
  maxWidth: 'xs' | 'sm' | 'md' | 'lg' | 'xl';
  title: string;
  children: React.ReactNode;
};

export default function Dialogue(props: DialogueProps) {
  // const [scroll, setScroll] = React.useState<'paper'>('paper');

  // const descriptionElementRef = React.useRef(null);
  // React.useEffect(() => {
  //   if (props.open) {
  //     const { current: descriptionElement } = descriptionElementRef;
  //     if (descriptionElement !== null) {
  //       descriptionElement.focus();
  //     }
  //   }
  // }, [props.open]);

  return (
    <div>
      <Dialog
        fullScreen={!!props.fullScreen}
        open={props.open}
        onClose={props.handleClose}
        // scroll={scroll}
        maxWidth={props.maxWidth}
        fullWidth
        aria-labelledby="scroll-dialog-title"
        aria-describedby="scroll-dialog-description"
      >
        <DialogTitle>
          <Typography variant="h2">{props.title}</Typography>
        </DialogTitle>
        <DialogContent>{props.children}</DialogContent>
        <DialogActions>
          <Button onClick={props.handleClose} color="primary">
            Close
          </Button>
        </DialogActions>
      </Dialog>
    </div>
  );
}
