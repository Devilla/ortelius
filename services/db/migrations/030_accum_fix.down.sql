
drop index `output_addresses_accumulate_in_output_processed` on `output_addresses_accumulate_in`;
drop index `output_addresses_accumulate_out_output_processed` on `output_addresses_accumulate_out`;

drop index `output_addresses_accumulate_in_output_id` on `output_addresses_accumulate_in`;
drop index `output_addresses_accumulate_out_output_id` on `output_addresses_accumulate_out`;

alter table `output_addresses_accumulate_in` drop COLUMN `output_in_processed`;
alter table `output_addresses_accumulate_in` drop COLUMN `output_out_processed`;
alter table `output_addresses_accumulate_out` drop COLUMN `output_out_processed`;
