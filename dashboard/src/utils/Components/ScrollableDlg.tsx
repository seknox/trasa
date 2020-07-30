import React from 'react';
import Button from '@material-ui/core/Button';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogContentText from '@material-ui/core/DialogContentText';
import DialogTitle from '@material-ui/core/DialogTitle';
import Typography from '@material-ui/core/Typography';

type DialogueProps = {
  open: boolean;
  handleClose: () => void;
  maxWidth: 'xs' | 'sm' | 'md' | 'lg' | 'xl';
  title: string;
  children: React.ReactNode;
  data: any;
};

export default function ScrollDialog(props: DialogueProps) {
  // const [open, setOpen] = React.useState(false);
  // const [scroll, setScroll] = React.useState<'paper'>('paper');

  // const handleClickOpen = () => {
  //   setOpen(true);
  //   setScroll('paper');
  // };

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
        open={props.open}
        onClose={props.handleClose}
        scroll="paper"
        maxWidth="lg"
        fullWidth
        aria-labelledby="scroll-dialog-title"
        aria-describedby="scroll-dialog-description"
      >
        <DialogTitle>
          <Typography variant="h2">{props.title}</Typography>
        </DialogTitle>
        <DialogContent dividers>
          <DialogContentText
            id="scroll-dialog-description"
            // ref={descriptionElementRef}
            tabIndex={-1}
          >
            {props.data}
          </DialogContentText>
          {props.children}
        </DialogContent>
        <DialogActions>
          <Button onClick={props.handleClose} color="primary">
            Close
          </Button>
        </DialogActions>
      </Dialog>
    </div>
  );
}
