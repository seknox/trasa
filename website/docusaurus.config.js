module.exports = {
  title: 'TRASA by Seknox',
  tagline: 'Manage, control and monitor access to your internal infrastructure',
  url: 'https://www.trasa.io',
  baseUrl: '/',
  favicon: 'img/trasa-icon.png',
  organizationName: 'Seknox',
  projectName: 'TRASA',
  onBrokenLinks: 'ignore',
  // plugins: ['@docusaurus/plugin-google-analytics'],
  themeConfig: {
    colorMode: {
      disableSwitch: true,
    },
    // sidebarCollapsible: false,
    googleAnalytics: {
      trackingID: 'UA-126179756-2',
      // Optional fields.
      anonymizeIP: false, // Should IPs be anonymized?
    },
    navbar: {
      title: '',
      logo: {
        alt: 'TRASA Logo',
        src: 'img/trasa.svg',
      },
      items: [
        // {
        //   label: 'Use Cases',
        //   position: 'left',
        //   items: [
        //     {
        //       to: 'use-cases/zero-trust-service-access',
        //       label: 'Zero Trust Service Access',
        //       position: 'right',
        //     },
        //     {
        //       to: 'use-cases/device-trust',
        //       label: 'Device Trust',
        //       position: 'right',
        //     },
        //     {
        //       to: 'use-cases/two-factor-authentication',
        //       label: 'Two Factor Authentication',
        //       position: 'right',
        //     },
        //     {
        //       to: 'use-cases/privileged-access-management',
        //       label: 'Privileged Access Management',
        //       position: 'right',
        //     },
        //   ],
        // },

        {
          label: 'Community',
          position: 'right',
          items: [
            { to: 'blog', label: 'Blog', position: 'right' },
            {
              to:
                'https://github.com/seknox/trasa/issues?q=is%3Aissue+is%3Aopen+label%3A%22good+first+issue%22',
              label: 'Contribute',
              position: 'right',
            },
            { to: 'https://discuss.trasa.io', label: 'Forum', position: 'right' },
          ],
        },
        // { to: 'features', label: 'Features', position: 'left' },
        // { to: 'security', label: 'Security', position: 'left' },
        {
          to: 'docs/',
          activeBasePath: 'docs',
          label: 'Admin Docs',
          position: 'right',
        },

        {
          to: 'docs/guides/getting-started',
          activeBasePath: 'guides',
          label: 'User Guides',
          position: 'right',
        },

        {
          href: 'https://github.com/seknox/trasa',
          className: 'header-github-link',
          position: 'right',
          'aria-label': 'GitHub repository',
          // label: 'Github',
        },
        {
          href: 'https://www.trasa.io/docs/install/installation/',
          className: 'signup-link',
          position: 'right',
          label: 'Try Now',
          'aria-label': 'GitHub repository',
        },
      ],
    },
    // announcementBar: {
    //   id: 'support_us', // Any value that will identify this message.
    //   content:
    //     'We are looking to revamp our docs, please fill <a target="_blank" rel="noopener noreferrer" href="#">this survey</a>',
    //   backgroundColor: '#fafbfc', // Defaults to `#fff`.
    //   textColor: '#091E42',
    // },
    footer: {
      style: 'dark',
      links: [
        // {
        //   title: 'Community',
        //   items: [
        //     {
        //       label: 'Stack Overflow',
        //       href: 'https://stackoverflow.com/questions/tagged/docusaurus',
        //     },
        //     {
        //       label: 'Discord',
        //       href: 'https://discordapp.com/invite/docusaurus',
        //     },
        //     {
        //       label: 'Twitter',
        //       href: 'https://twitter.com/docusaurus',
        //     },
        //   ],
        // },
        // {
        //   title: 'More',
        //   items: [
        //     {
        //       label: 'Blog',
        //       to: 'blog',
        //     },
        //     {
        //       label: 'GitHub',
        //       href: 'https://github.com/facebook/docusaurus',
        //     },
        //   ],
        // },
      ],
      copyright: `Copyright © ${new Date().getFullYear()} TRASA by Seknox Pte. Ltd.`,
    },
  },
  presets: [
    [
      '@docusaurus/preset-classic',
      {
        docs: {
          // It is recommended to set document id as docs home page (`docs/` path).
          // homePageId: 'getting-started/overview',
          sidebarPath: require.resolve('./sidebars.js'),
          // Please change this to your repo.
          editUrl: 'https://github.com/seknox/trasa/edit/master/website/',
        },
        blog: {
          showReadingTime: true,
          // Please change this to your repo.
          editUrl: 'https://github.com/seknox/trasa/edit/master/website/blog/',
        },
        theme: {
          customCss: require.resolve('./src/css/custom.css'),
        },
      },
    ],
  ],
};
