// Libraries
import React, {PureComponent, ChangeEvent} from 'react'
import {getDeep} from 'src/utils/wrappers'

// Components
import {ErrorHandling} from 'src/shared/decorators/errors'
import {
  Button,
  ComponentColor,
  ComponentSize,
  Input,
  InputType,
  Form,
  Columns,
  IconFont,
  Grid,
  ComponentStatus,
} from 'src/clockface'

// APIS
import {setSetupParams, SetupParams, signin} from 'src/onboarding/apis'

// Constants
import * as copy from 'src/shared/copy/notifications'

// Types
import {StepStatus} from 'src/clockface/constants/wizard'
import {OnboardingStepProps} from 'src/onboarding/containers/OnboardingWizard'

interface State extends SetupParams {
  confirmPassword: string
  isAlreadySet: boolean
  isPassMismatched: boolean
}

@ErrorHandling
class AdminStep extends PureComponent<OnboardingStepProps, State> {
  constructor(props: OnboardingStepProps) {
    super(props)
    const {setupParams} = props
    this.state = {
      username: getDeep(setupParams, 'username', ''),
      password: getDeep(setupParams, 'password', ''),
      confirmPassword: getDeep(setupParams, 'password', ''),
      org: getDeep(setupParams, 'org', ''),
      bucket: getDeep(setupParams, 'bucket', ''),
      isAlreadySet: !!setupParams,
      isPassMismatched: false,
    }
  }

  public componentDidMount() {
    window.addEventListener('keydown', this.handleKeydown)
  }

  public componentWillUnmount() {
    window.removeEventListener('keydown', this.handleKeydown)
  }

  public render() {
    const {username, password, confirmPassword, org, bucket} = this.state
    const icon = this.InputIcon
    const status = this.InputStatus
    return (
      <div className="onboarding-step">
        <h3 className="wizard-step--title">Setup Admin User</h3>
        <h5 className="wizard-step--sub-title">
          You will be able to create additional Users, Buckets and Organizations
          later
        </h5>
        <Form className="onboarding--admin-user-form">
          <Grid>
            <Grid.Row>
              <Grid.Column widthXS={Columns.Ten} offsetXS={Columns.One}>
                <Form.Element label="Admin Username">
                  <Input
                    value={username}
                    onChange={this.handleUsername}
                    titleText="Admin Username"
                    size={ComponentSize.Medium}
                    icon={icon}
                    status={status}
                    disabledTitleText="Admin username has been set"
                    autoFocus={true}
                  />
                </Form.Element>
              </Grid.Column>
              <Grid.Column widthXS={Columns.Five} offsetXS={Columns.One}>
                <Form.Element label="Admin Password">
                  <Input
                    type={InputType.Password}
                    value={password}
                    onChange={this.handlePassword}
                    titleText="Admin Password"
                    size={ComponentSize.Medium}
                    icon={icon}
                    status={status}
                    disabledTitleText="Admin password has been set"
                  />
                </Form.Element>
              </Grid.Column>
              <Grid.Column widthXS={Columns.Five}>
                <Form.Element label="Confirm Admin Password">
                  <Input
                    type={InputType.Password}
                    value={confirmPassword}
                    onChange={this.handleConfirmPassword}
                    titleText="Confirm Admin Password"
                    size={ComponentSize.Medium}
                    icon={icon}
                    status={this.passwordStatus}
                    disabledTitleText="Admin password has been set"
                  />
                </Form.Element>
              </Grid.Column>
              <Grid.Column widthXS={Columns.Ten} offsetXS={Columns.One}>
                <Form.Element label="Default Organization Name">
                  <Input
                    value={org}
                    onChange={this.handleOrg}
                    titleText="Default Organization Name"
                    size={ComponentSize.Medium}
                    icon={icon}
                    status={ComponentStatus.Default}
                    placeholder="Your organization is where everything you create lives"
                    disabledTitleText="Default organization name has been set"
                  />
                </Form.Element>
              </Grid.Column>
              <Grid.Column widthXS={Columns.Ten} offsetXS={Columns.One}>
                <Form.Element label="Default Bucket Name">
                  <Input
                    value={bucket}
                    onChange={this.handleBucket}
                    titleText="Default Bucket Name"
                    size={ComponentSize.Medium}
                    icon={icon}
                    status={status}
                    placeholder="Your bucket is where you will store all your data"
                    disabledTitleText="Default bucket name has been set"
                  />
                </Form.Element>
              </Grid.Column>
            </Grid.Row>
          </Grid>
          <div className="wizard--button-container">
            <div className="wizard--button-bar">
              <Button
                color={ComponentColor.Default}
                text="Back to Start"
                size={ComponentSize.Medium}
                onClick={this.props.onDecrementCurrentStepIndex}
              />
              <Button
                color={ComponentColor.Primary}
                text="Continue to Data Loading"
                size={ComponentSize.Medium}
                onClick={this.handleNext}
                status={this.nextButtonStatus}
                titleText={this.nextButtonTitle}
              />
            </div>
          </div>
        </Form>
      </div>
    )
  }

  private handleUsername = (e: ChangeEvent<HTMLInputElement>): void => {
    const username = e.target.value
    this.setState({username})
  }

  private handlePassword = (e: ChangeEvent<HTMLInputElement>): void => {
    const {confirmPassword} = this.state
    const password = e.target.value
    const isPassMismatched = confirmPassword && password !== confirmPassword
    this.setState({password, isPassMismatched})
  }

  private handleConfirmPassword = (e: ChangeEvent<HTMLInputElement>): void => {
    const {password} = this.state
    const confirmPassword = e.target.value
    const isPassMismatched = confirmPassword && password !== confirmPassword
    this.setState({confirmPassword, isPassMismatched})
  }

  private handleOrg = (e: ChangeEvent<HTMLInputElement>): void => {
    const org = e.target.value
    this.setState({org})
  }

  private handleBucket = (e: ChangeEvent<HTMLInputElement>): void => {
    const bucket = e.target.value
    this.setState({bucket})
  }

  private handleKeydown = (e: KeyboardEvent): void => {
    if (
      e.key === 'Enter' &&
      this.nextButtonStatus === ComponentStatus.Default
    ) {
      this.handleNext()
    }
  }

  private get nextButtonStatus(): ComponentStatus {
    const {
      username,
      password,
      org,
      bucket,
      confirmPassword,
      isPassMismatched,
    } = this.state
    if (
      username &&
      password &&
      confirmPassword &&
      org &&
      bucket &&
      !isPassMismatched
    ) {
      return ComponentStatus.Default
    }
    return ComponentStatus.Disabled
  }

  private get nextButtonTitle(): string {
    const {username, password, org, bucket, confirmPassword} = this.state
    if (username && password && confirmPassword && org && bucket) {
      return 'Next'
    }
    return 'All fields are required to continue'
  }

  private get passwordStatus(): ComponentStatus {
    const {isAlreadySet, isPassMismatched} = this.state
    if (isAlreadySet) {
      return ComponentStatus.Disabled
    }
    if (isPassMismatched) {
      return ComponentStatus.Error
    }
    return ComponentStatus.Default
  }

  private get InputStatus(): ComponentStatus {
    const {isAlreadySet} = this.state
    if (isAlreadySet) {
      return ComponentStatus.Disabled
    }
    return ComponentStatus.Default
  }

  private get InputIcon(): IconFont {
    const {isAlreadySet} = this.state
    if (isAlreadySet) {
      return IconFont.Checkmark
    }
    return null
  }

  private handleNext = async () => {
    const {
      onSetStepStatus,
      currentStepIndex,
      handleSetSetupParams,
      notify,
      onIncrementCurrentStepIndex,
    } = this.props

    const {username, password, org, bucket, isAlreadySet} = this.state

    if (isAlreadySet) {
      onSetStepStatus(currentStepIndex, StepStatus.Complete)
      onIncrementCurrentStepIndex()
      return
    }

    const setupParams = {
      username,
      password,
      org,
      bucket,
    }

    try {
      await setSetupParams(setupParams)
      await signin({username, password})
      notify(copy.SetupSuccess)
      handleSetSetupParams(setupParams)
      onSetStepStatus(currentStepIndex, StepStatus.Complete)
      onIncrementCurrentStepIndex()
    } catch (error) {
      notify(copy.SetupError)
    }
  }
}

export default AdminStep
