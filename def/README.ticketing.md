# Ticketing The Agony and the Ecstasy of Event planning.

### Capacity and stock

The entry point for thi model are `EventSlot` items, one of these specifies Slot of Event time for which we ask an admitance ticket, these can overlap, for instance **General Admitance** Will be required to enter the general meeting area of one Event and it is a precondition to be able to obtain a **Free Workshop Ticket** which grants participation to a worhkshop that happens in the same time frame as the general event for which we granted admittance.

`EventSlot`s are not necesarily a representation of a full Slot, as an example, for a **General Admitance** that allows 200 people due to venue contraints, the Slots can be as follows:

(This example is not representative at all on how events have typically been handled, it is just a way to stretch the model to the corners)

 * 50 Early Bird Tickets, Cost 2U$D, Sold from January to March, available to public and valid to redeem from 1st to last day of the conference
 * 100 Regular Tickets, Cost 4U$D, Sold from March to August, available to public and valid to redeem from 1st to last day of the conference
 * 50 Sponsor Tickets, Cost 0U$D, Sold from January to August, not available to public and valid to redeem from 1st to last day of the conference

### Actors

Our main actor is the `Attendee` which is a human being (or entity? entity rep?) that will buy and optionally use these Slot tickets (they could buy and transfer or just be a sponsor representative and transfer the tickets to their employees assisting the `Event`)

An actor claims slots through `SlotClaim`s, each of those have an unique Ticket ID (the ticket number but it's not a number), a claim can be tranfered while:

 * It is not redeemed.
 * The `SlotClaim.EventSlot.StartTime` is in the future.

**Constraint** The total amount of `SlotClaim` cannot be more than `SlotClaim.EventSlot.Capacity`

Redeeming a ticket requires 3 piece of Information:
  
 * The attendee assisting to event email.
 * Said attendee having accepted a Code of Conduct.
 * The ticket (or tickets) claimed by the attendee.

### The vile metal

Events cost money and so we need to charge admittance to them.

As you have already read in the previous section, each "purchased" admittance is a `SlotClaim` and as such, these are the items to be paid. You can think of them as `EventSlot` being a product, lets say *Milk In Boxes* and `SlotClaim` being an unit of it *As carton of milk* (or whatever container milks comes in where you live).

Every `SlotClaim` must be backed by a `ClaimPayment` bear in mind that a `ClaimPayment` might be paying for multiple `SlotClaim` (ie, I buy GA and Workshop and GophersWhoDontGluten party and I will pay for all of these). A `ClaimPayment` has a reference to N `SlotClaim` which are being payed. The `TotalDue()` is therefore the sum of `SlotClaim.EventSlot.Cost`. It also has a list of payments which are made through `FinancialInstrument` which is anything that can fulfill the role of money. 
The available **types** of Financial Instruments are:
 
 * Cash: Money I have "In hand" (or to be precise, in bank account).
 * Discount: Money We have not collected nor we intend to, it is written off as loss (this is in accounting terms, not meant as a negative thing, assuming you consider you total patrimony be X*10 where X is the value which you intend to obtain from selling an Item an 10 your total stock of items, it means that you intended to earn X for an item sell and now you earn X-Y therefore reducing your total patrimony, it is far more complex than this but the simplification suffices I think)
 * Receivables: Money we have not collected but we intendo to, in patrimonial terms, this is money you own but not yet have in hand so it is yours but you can't use it yet (much like a videogame you purchased, it is yours but you cannot use it until it finihes downloading 50G of updates :p).

Currently we have the following `FinancialInstrument`:

 * Money (`PaymentMethodMoney`) this is.. well, money. It holds an amount and a ref (in whatever format it is handed to us) to a payment in stripe.
 * Credit (`PaymentMethodCreditNote`) this is, in few words, a debt from the attendee to the event. Typically this is used when invoicing needs to happen before the payment i processed (like with large corporation sponsors).
 * Discount (`PaymentMethodConferenceDiscount`) the event organizer can discount the price of an item by giving someone scholarships or benefits or whatever kind of discount in a `SlotClaim`

#### Balances

A couple of simplifications were added as sample operations to our use of money.

**Note:** I was quite liberal with the use of accounting terms in this section an the code backing it.

* `DebtBalanced`: A debt is balanced when we have at least as many cash+discounts as Receivables (if credit was never issued then it is simply that we have either positive balance or no debt)
* `PaymentBalanced`: The total of credit plus cash plus discount covers the passed amount (typically used to determine if it covers a ticket and therefore it can be invoiced.)

A `ClaimPayment` is `Fulfulled` when `DebtBalanced` and `PaymentBalanced` of `ClaimPayment.Total()` are true.

#### Invoices

Invoces will be generated by... something, which we do not know, if it has an API and returns references, all the best.

An invoice covers one `ClaimPayment` its items are `ClaimPayment.ClaimsPayed`, each line item cost is `ClaimsPayment.ClaimsPayed[i].EventSlot.Cost` and its total is `ClaimPayment.Total()`. Ideally we then fill `ClaimPayment.Invoice` with the ref.

If the US has the concept of a "Mark as payed" seal in invoices, `ClaimPayment.Fulfilled()` returning true is the indicator that it should be applied.

