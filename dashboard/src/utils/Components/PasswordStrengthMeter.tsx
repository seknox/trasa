import React from 'react';
import zxcvbn from 'zxcvbn';
import './PassStrengthMeter.css';

function PasswordStrengthMeter(props: any) {
  const createPasswordLabel = (result: any) => {
    const sc = result.score; // - 2
    switch (sc) {
      case 0:
        return 'Weak';
      case 1:
        return 'Weak';
      case 2:
        return 'Fair';
      case 3:
        return 'Good';
      case 4:
        return 'Strong';
      default:
        return 'Weak';
    }
  };

  const { password } = props;
  const testedResult = zxcvbn(password);

  return (
    <div className="password-strength-meter">
      <progress
        className={`password-strength-meter-progress strength-${createPasswordLabel(testedResult)}`}
        value={testedResult.score}
        max="5"
      />
  

        <label className="password-strength-meter-label">
         
          {password && (
            <div>
              <strong>Password strength:</strong> {createPasswordLabel(testedResult)}
            </div>
          )}
        </label>

    </div>
  );
}

export default PasswordStrengthMeter;
