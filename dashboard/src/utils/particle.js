import { makeStyles } from '@material-ui/core/styles';
import React from 'react';
import Particles from 'react-particles-js';

// STYLES
const useStyles = makeStyles((theme) => ({
  root: {
    background: 'rgba(1,1,35,1)', //'#0b1728ff', //'#0b1728ff',
    width: '100%',
    Height: '100%',
    minHeight: 900,
    margin: '0'
  },
  login: {
    position: 'absolute',
    top: '45%',
    right: '50%',
    transform: 'translate(50%,-50%)'
  },
  tfa: {
    position: 'absolute',
    top: '45%',
    right: '50%',
    transform: 'translate(50%,-50%)'
  },
  particle: {
    height: '100%',
    minHeight: 900,
    background: 'rgba(1,1,35,1)'
  }
}));

export default function Part(props) {
  const classes = useStyles();
  return (
    <div className={classes.particle}>
      <Particles
        params={{
          particles: {
            number: {
              value: 100,
              density: {
                enable: false,
                value_area: 1000
              }
            },
            color: {
              value: '#ffffff'
            },
            shape: {
              type: 'square',
              stroke: {
                width: 2,
                color: '#000000'
              },
              polygon: {
                nb_sides: 3
              },
              image: {
                src: 'img/github.svg',
                width: 100,
                height: 100
              }
            },
            opacity: {
              value: 0.5,
              random: false,
              anim: {
                enable: false,
                speed: 1,
                opacity_min: 0.1,
                sync: false
              }
            },
            size: {
              value: 3,
              random: true,
              anim: {
                enable: false,
                speed: 40,
                size_min: 0.1,
                sync: false
              }
            },
            line_linked: {
              enable: true,
              distance: 150,
              color: '#0eafb9', //"#0014ff",
              opacity: 0.4,
              width: 1.5
            },
            move: {
              enable: true,
              speed: 1,
              direction: 'none',
              random: true,
              straight: false,
              out_mode: 'bounce',
              bounce: false,
              attract: {
                enable: true,
                rotateX: 600,
                rotateY: 1200
              }
            }
          },
          interactivity: {
            detect_on: 'canvas',
            events: {
              onhover: {
                enable: true,
                mode: 'grab'
              },
              onclick: {
                enable: true,
                mode: 'push'
              },
              resize: true
            },
            modes: {
              grab: {
                distance: 140,
                line_linked: {
                  opacity: 1
                }
              },
              bubble: {
                distance: 400,
                size: 40,
                duration: 2,
                opacity: 8,
                speed: 3
              },
              repulse: {
                distance: 200,
                duration: 0.4
              },
              push: {
                particles_nb: 4
              },
              remove: {
                particles_nb: 2
              }
            }
          },
          retina_detect: true
        }}
      />
    </div>
  );
}
