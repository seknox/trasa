module.exports = {
  docs: [
    {
      'Getting Started': [
        'getting-started/overview',
        'getting-started/concepts',
        'getting-started/how-to',
        'getting-started/glossary',
        // 'getting-started/signup-or-install',
      ],
      Install: ['install/installation', 'install/initial-setup'],
      Tutorial: [
        'tutorial/intro',
        'tutorial/setup-trasa-server',
        'tutorial/create-users',
        'tutorial/create-policy',
        'tutorial/protect-services',
        'tutorial/test-access-to-services',
        'tutorial/monitor-access',
        'tutorial/advance-configurations',
      ],
      'How to Access ?': [
        'how-to-access/dashboard-login',
        'how-to-access/tfa',
        'how-to-access/ssh-connection-via-proxy',
        'how-to-access/rdp-connection-via-proxy',
        'how-to-access/https-connection-via-proxy',
      ],
      Providers: [
        'providers/providers',
        { users: ['providers/users/ldap/ldap'], vault: ['providers/vault/tsxvault'] },
      ],
      Users: [
        'users/users',
        'users/creating-updating-users',
        'users/account-setup',
        'users/change-pass',
        // 'users/trasaIDP/trasa-idp',
        // 'users/ldap/ldap',
        // 'users/ad/active-directory',
        // 'users/freeIPA/free-ipa',
      ],
      'Device Trust': ['device-trust/enrol-2fa-device', 'device-trust/install-trasa-device-agent'],
      Services: [
        'services/introduction',
        'services/privilege',
        'services/http/http-service',
        'services/rdp/rdp-service',
        'services/ssh/ssh-service',
        'services/radius/radius-server',
      ],
      Policies: [
        'policies/policy-introduction',
        'policies/basic-policy',
        'policies/device-policy',
        'policies/adhoc-policy',
      ],
      'Access Map': [
        'access-map/access-mapping',
        'access-map/dynamic-access',
        'access-map/adhoc-access',
      ],
      'Access Proxy': ['access-proxy/introduction', 'access-proxy/proxy-vs-direct'],
      'Two Factor Authentication': [
        'native-tfa/two-factor-authentication',
        'native-tfa/windows/windows-two-factor-authentication',
        'native-tfa/linux-two-factor-authentication',
      ],
      // 'Secret Vault': ['providers/secret-vault/index'],
      // 'Cloud Providers': ['cloud/amazon-web-services', 'cloud/google-cloud', 'cloud/digital-ocean'],

      'System Configurations': [
        'system/fcm-setup',
        'system/email-setup',
        'system/config-reference',
      ],
    },
  ],
  // guides: [
  //   'guides/getting-started',
  //   {
  //     Account: ['guides/user/account/password-setup'],
  //     Device: [
  //       'guides/user/device/enrol-2fa-device',
  //       'guides/user/device/install-trasa-device-agent',
  //     ],
  //     Access: [
  //       //  'guides/user/access/dashboard-login',
  //       'guides/user/access/adhoc-access',
  //       'guides/user/access/tfa',
  //       'guides/user/access/ssh-connection-via-proxy',
  //       'guides/user/access/rdp-connection-via-proxy',
  //       'guides/user/access/https-connection-via-proxy',
  //     ],
  //   },
  // ],
};
