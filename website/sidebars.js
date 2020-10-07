module.exports = {
  docs: [
    {
      'Getting Started': [
        'getting-started/overview',
        'getting-started/concepts',
        'getting-started/glossary',
        'getting-started/signup-or-install',

        'getting-started/how-to',
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
      Users: [
        'users/users',
        'users/creating-updating-users',
        // 'users/trasaIDP/trasa-idp',
        // 'users/ldap/ldap',
        // 'users/ad/active-directory',
        // 'users/freeIPA/free-ipa',
      ],
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
      'Access Map': ['access-map/access-mapping', 'access-map/dynamic-access'],
      'Access Proxy': ['access-proxy/introduction'],
      'Two Factor Authentication': [
        'native-tfa/two-factor-authentication',
        'native-tfa/windows/windows-two-factor-authentication',
        'native-tfa/linux-two-factor-authentication',
      ],
      'Secret Vault': ['providers/secret-vault/index'],
      // 'Cloud Providers': ['cloud/amazon-web-services', 'cloud/google-cloud', 'cloud/digital-ocean'],
      Providers: [{ users: ['providers/users/ldap/ldap'] }],

      'System Configurations': [
        'system/fcm-setup',
        'system/email-setup',
        'system/config-reference',
      ],
    },
  ],
  guides: [
    'guides/getting-started',
    {
      Account: ['guides/user/account/password-setup'],
      Device: [
        'guides/user/device/enrol-2fa-device',
        'guides/user/device/install-trasa-device-agent',
      ],
      Access: [
        'guides/user/access/dashboard-login',
        'guides/user/access/adhoc-access',
        'guides/user/access/tfa',
        'guides/user/access/ssh-connection-via-proxy',
        'guides/user/access/rdp-connection-via-proxy',
        'guides/user/access/https-connection-via-proxy',
      ],
    },
  ],
};
