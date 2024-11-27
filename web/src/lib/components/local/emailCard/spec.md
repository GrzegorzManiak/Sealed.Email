# Email card component

This component is used to display a little bit of information about an 
email. It is heavily inspired by the Outlook email card, although it is
not a direct copy.

## Functionality

### (Default) No attachments, No chain

The email card should display the following information:

- Sender's name OR email address if the name is not available
- Subject of the email
- Date and time the email was received
- The first line of the email body
- A 'Favorite' button that can be toggled on and off
- A 'Untrusted' icon if the email is from an untrusted sender
- A 'Pin' icon to allow the user to pin the email to the top of the list
- A 'Delete' icon to allow the user to delete the email
- A 'Read' 'Unread' toggle button to allow the user to mark the email as read or unread

### Attachments

If the email has attachments, the email card should 
display all of the above information, as well as:

- An extension of the bottom of the email card that shows the number of attachments
- A 'Showcase' attachment, which is the first attachment in the list, and is displayed as a small icon with the file name underneath it
- The icon should math the file type of the attachment, else just a file icon should be displayed
- A '+X' text to the right of the button, where X is the number of attachments that are not displayed in the showcase

### Chain

If the email is part of a chain, the email card should display all of the default information, as well as:

- A more compact version of the email card, with the sender's name and the first line of the email body, no subject.
- No profile picture should be displayed, nor a favorite button
- No 'pin', only 'delete' and 'read/unread' buttons
- A tag indicating if the user is the sender of the email, eg 'You'
- In the closed state, an extension to the bottom of the email card that shows the number of emails in the chain and a 'Show' button

## Selection Mode

Emails can be selected in groups, when hovering over the email card a check box should appear where the
profile picture would be. When the email is selected, all of the email cards should shift into this mode
where the background of the email cards turns gray, and they all display a check box in the same location.

When an individual email is clicked and it is not in the selection mode, the email card should turn blue.
If the email has a chain, all emails in the chain should turn blue, and if any email within the chain is clicked
they should turn to a different shade of blue to indicate that they are selected.